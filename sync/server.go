package sync

import (
	"fmt"
	"sync"

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
	conns      map[string]Context
	sync.RWMutex
}

func (v *server) Start() error {
	return nil
}

func (v *server) Stop() error {
	return nil
}

func (v *server) Publish(p Packet) error {
	b, err := p.Bytes()
	if err != nil {
		return err
	}
	if err := v.message.Emit(v.topic, b); err != nil {
		return err
	}
	return nil
}

func (v *server) Subscribe(p Packet) error {
	return nil
}

func (v *server) Handle(r router.Context, n *websocket.Conn) {
	u := fmt.Sprintf("%s-%s", r.MachineID(), r.ProcessID())
	c := &connection{
		machineID: r.MachineID(),
		processID: r.ProcessID(),
		ctx:       r.Context(),
		channels:  make([]Channel, 0),
		req:       r.Req(),
		res:       r.Res(),
		client:    r.Client(),
		cookie:    r.Cookie(),
		header:    r.Header(),
		cache:     v.cache,
		message:   v.message,
		server:    v,
		conn:      n,
	}
	v.RLock()
	v.conns[u] = c
	v.RUnlock()
	c.BeforeOpen()
	defer func() {
		v.RLock()
		delete(v.conns, u)
		v.RUnlock()
		c.AfterClose()
	}()
	for {
		mt, message, err := n.ReadMessage()
		if err != nil {
			break
		}
		v.Publish(&packet{
			T: mt,
			M: message,
		})
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
