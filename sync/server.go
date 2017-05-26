package sync

import (
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/samuelngs/hyper/cache"
	"github.com/samuelngs/hyper/message"
	"github.com/samuelngs/hyper/router"
)

type server struct {
	id         string
	topic      []byte
	cache      cache.Service
	message    message.Service
	namespaces []Namespace
	nsmap      map[string]Namespace
	conns      map[string]Context
	hookbo     HookFunc
	hookac     HookFunc
	stop       message.Close
	sync.RWMutex
}

func (v *server) Start() error {
	v.Lock()
	for _, n := range v.namespaces {
		c := n.Config()
		s := c.Namespace()
		a := c.Aliases()
		v.nsmap[s] = n
		for _, i := range a {
			v.nsmap[i] = n
		}
	}
	v.Unlock()
	go func() {
		v.stop = v.message.Listen(v.topic, func(b []byte) {
			d := &Distribution{}
			if err := proto.Unmarshal(b, d); err != nil {
				return
			}
			v.Subscribe(d)
		})
	}()
	return nil
}

func (v *server) Stop() error {
	v.Lock()
	v.nsmap = make(map[string]Namespace)
	v.Unlock()
	v.stop()
	return nil
}

func (v *server) Publish(d *Distribution) error {
	b, err := proto.Marshal(d)
	if err != nil {
		return err
	}
	return v.message.Emit(v.topic, b)
}

func (v *server) Subscribe(d *Distribution) error {
	if d.Packet == nil {
		return InvalidPacket.Fill()
	}
	n := v.Namespace(d.Packet.GetNamespace())
	cs := n.Channels()
	if !cs.Has(d.Packet.GetChannel()) {
		return ChannelNotExist.Fill(d.Packet.GetChannel())
	}
	ch := cs.Get(d.Packet.GetChannel())
	ns := ch.NodeSubscribers()
	for _, c := range ns {
		c.Write(d.Packet)
	}
	return nil
}

func (v *server) BeforeOpen(f HookFunc) {
	v.hookbo = f
}

func (v *server) AfterClose(f HookFunc) {
	v.hookac = f
}

func (v *server) Handle(r router.Context, n *websocket.Conn) {
	u := fmt.Sprintf("%s-%s", r.MachineID(), r.ProcessID())
	c := &connection{
		machineID:     r.MachineID(),
		processID:     r.ProcessID(),
		identity:      new(identity),
		subscriptions: &subscriptions{make([]Channel, 0)},
		ctx:           r.Context(),
		req:           r.Req(),
		res:           r.Res(),
		client:        r.Client(),
		cookie:        r.Cookie(),
		header:        r.Header(),
		cache:         v.cache,
		message:       v.message,
		server:        v,
		conn:          n,
	}
	v.Lock()
	v.conns[u] = c
	v.Unlock()
	if v.hookbo != nil {
		v.hookbo(c)
	}
	c.BeforeOpen()
	defer func() {
		recover()
		v.Lock()
		delete(v.conns, u)
		v.Unlock()
		if v.hookac != nil {
			v.hookac(c)
		}
		c.AfterClose()
	}()
	for {
		mt, message, err := n.ReadMessage()
		if err != nil {
			break
		}
		v.Read(mt, message, c)
	}
}

func (v *server) Read(mt int, message []byte, c Context) {
	if mt == websocket.BinaryMessage && message != nil && len(message) > 0 {
		p := &Packet{}
		// parse protobuf packet
		if err := proto.Unmarshal(message, p); err != nil {
			c.Write(&Packet{
				Action: ActionMessageFailure,
				Error:  InvalidPacket.Fill().JsonString(),
			})
			return
		}
		v.RLock()
		n, ok := v.nsmap[p.GetNamespace()]
		v.RUnlock()
		// namespace does not exist, send error
		if !ok {
			c.Write(&Packet{
				Action: ActionMessageFailure,
				Error:  NamespaceNotExist.Fill(p.GetNamespace()).JsonString(),
			})
			return
		}
		switch p.GetAction() {
		case ActionSubscribe:
			v.HandleSubscribe(p, n, c)
		case ActionUnsubscribe:
			v.HandleUnsubscribe(p, n, c)
		case ActionMessage:
			v.HandleMessage(p, n, c)
		default:
			c.Write(&Packet{
				Action: ActionMessageFailure,
				Error:  InvalidAction.Fill(p.GetAction()).JsonString(),
			})
		}
	}
}

func (v *server) HandleSubscribe(p *Packet, n Namespace, c Context) {
	// checks if channel is already subscribed
	if subscribed := c.Subscriptions().Has(p.GetNamespace(), p.GetChannel()); subscribed {
		c.Write(&Packet{
			Action: ActionSubscribeFailure,
			Error:  ChannelAlreadySubscribed.Fill(p.GetNamespace(), p.GetChannel()).JsonString(),
		})
		return
	}
	if fn := n.Config().Authorize(); fn != nil {
		if err := fn(p.GetChannel(), c); err != nil {
			c.Write(&Packet{
				ID:        p.GetID(),
				Action:    ActionSubscribeFailure,
				Namespace: p.GetNamespace(),
				Channel:   p.GetChannel(),
				Error:     ChannelUnauthorized.Fill(p.GetNamespace(), p.GetChannel()).JsonString(),
			})
			return
		}
	}
	cs := n.Channels()
	if !cs.Has(p.GetChannel()) {
		cs.Add(p.GetChannel())
	}
	ch := cs.Get(p.GetChannel())
	if !ch.Has(c) {
		ch.Subscribe(c)
	}
	c.Write(&Packet{
		ID:        p.GetID(),
		Action:    ActionSubscribeSuccessful,
		Namespace: p.GetNamespace(),
		Channel:   p.GetChannel(),
	})
}

func (v *server) HandleUnsubscribe(p *Packet, n Namespace, c Context) {
	// checks if channel is not subscribed
	if subscribed := c.Subscriptions().Has(p.GetNamespace(), p.GetChannel()); !subscribed {
		c.Write(&Packet{
			Action: ActionUnsubscribeFailure,
			Error:  ChannelNotSubscribed.Fill(p.GetNamespace(), p.GetChannel()).JsonString(),
		})
		return
	}
	cs := n.Channels()
	if !cs.Has(p.GetChannel()) {
		return
	}
	ch := cs.Get(p.GetChannel())
	if ch.Has(c) {
		ch.Unsubscribe(c)
	}
	c.Write(&Packet{
		ID:        p.GetID(),
		Action:    ActionUnsubscribeSuccessful,
		Namespace: p.GetNamespace(),
		Channel:   p.GetChannel(),
	})
}

func (v *server) HandleMessage(p *Packet, n Namespace, c Context) {
	if subscribed := c.Subscriptions().Has(p.GetNamespace(), p.GetChannel()); !subscribed {
		c.Write(&Packet{
			Action: ActionMessageFailure,
			Error:  ChannelNotSubscribed.Fill(p.GetNamespace(), p.GetChannel()).JsonString(),
		})
		return
	}
	cs := n.Channels()
	ch := cs.Get(p.GetChannel())
	defer func() {
		if err := recover(); err != nil {
			if f := n.Config().Catch(); f != nil {
				f(p.GetMessage(), ch, c)
			}
		}
	}()
	if md := n.Config().Middlewares(); len(md) > 0 {
		for _, m := range md {
			m(p.GetMessage(), ch, c)
		}
	}
	if hs := n.Config().Handlers(); len(hs) > 0 {
		for _, h := range hs {
			if p.GetCall() == h.Name() {
				h.Func()(p.GetMessage(), ch, c)
			}
		}
	}
}

func (v *server) Namespace(s string) Namespace {
	for _, n := range v.namespaces {
		c := n.Config()
		if c.Namespace() == s {
			return n
		}
		for _, alias := range c.Aliases() {
			if alias == s {
				return n
			}
		}
	}
	n := &namespace{
		namespace: s,
		aliases:   make([]string, 0),
	}
	n.channels = &channels{
		namespace: n,
		channels:  make(map[string]Channel),
		server:    v,
	}
	v.namespaces = append(v.namespaces, n)
	return n
}

func (v *server) Namespaces() []Namespace {
	return v.namespaces
}

func (v *server) String() string {
	return "Hyper::Sync"
}
