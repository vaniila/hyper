package gws

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/vaniila/hyper/cache"
	"github.com/vaniila/hyper/message"
	"github.com/vaniila/hyper/router"
)

type server struct {
	id       string
	topic    []byte
	cache    cache.Service
	message  message.Service
	schema   graphql.Schema
	conns    map[string]Context
	tree     Store
	adaptor  router.GQLSubscriptionAdaptor
	hookauth AuthorizeFunc
	hookbo   HookFunc
	hookac   HookFunc
	stop     message.Close
	sync.RWMutex
}

func (v *server) Start() error {
	go func() {
		v.stop = v.message.Listen(v.topic, func(b []byte) {
			d := new(Distribution)
			if err := proto.Unmarshal(b, d); err != nil {
				return
			}
			v.Subscribe(d)
		})
	}()
	return nil
}

func (v *server) Stop() error {
	v.stop()
	return nil
}

func (v *server) Schema(schema graphql.Schema) {
	v.schema = schema
}

func (v *server) Adaptor() router.GQLSubscriptionAdaptor {
	return v.adaptor
}

func (v *server) Publish(d *Distribution) error {
	b, err := proto.Marshal(d)
	if err != nil {
		return err
	}
	return v.message.Emit(v.topic, b)
}

func (v *server) Subscribe(d *Distribution) error {
	if len(d.Field) == 0 || len(d.Payload) == 0 {
		return nil
	}
	var eqids, neids map[int64]struct{}
	var eqkeys, nekeys map[string]struct{}
	if d.Condition != nil {
		if l := len(d.Condition.EqIDs); l > 0 {
			eqids = make(map[int64]struct{}, l)
			for _, id := range d.Condition.EqIDs {
				eqids[id] = struct{}{}
			}
		}
		if l := len(d.Condition.NeIDs); l > 0 {
			neids = make(map[int64]struct{}, l)
			for _, id := range d.Condition.NeIDs {
				neids[id] = struct{}{}
			}
		}
		if l := len(d.Condition.EqKeys); l > 0 {
			eqkeys = make(map[string]struct{}, l)
			for _, key := range d.Condition.EqKeys {
				eqkeys[key] = struct{}{}
			}
		}
		if l := len(d.Condition.NeKeys); l > 0 {
			nekeys = make(map[string]struct{}, l)
			for _, key := range d.Condition.NeKeys {
				nekeys[key] = struct{}{}
			}
		}
	}
subscriptions:
	for _, sub := range v.tree.Get(d.Field) {
		if eqids != nil {
			if !sub.Connection().Identity().HasID() {
				continue
			}
			if _, ok := eqids[int64(sub.Connection().Identity().GetID())]; !ok {
				continue
			}
		}
		if neids != nil {
			if !sub.Connection().Identity().HasID() {
				continue
			}
			if _, ok := neids[int64(sub.Connection().Identity().GetID())]; ok {
				continue
			}
		}
		if eqkeys != nil {
			if !sub.Connection().Identity().HasKey() {
				continue
			}
			if _, ok := eqkeys[sub.Connection().Identity().GetKey()]; !ok {
				continue
			}
		}
		if nekeys != nil {
			if !sub.Connection().Identity().HasKey() {
				continue
			}
			if _, ok := nekeys[sub.Connection().Identity().GetKey()]; ok {
				continue
			}
		}
		args := sub.Arguments()
		for _, filter := range d.Filters {
			da, ok := args[filter.GetKey()]
			if !ok {
				continue subscriptions
			}
			switch o := filter.GetValOneof().(type) {
			case *Filter_StringValue:
				s, ok := da.(string)
				if !ok || s != o.StringValue {
					continue subscriptions
				}
			case *Filter_IntValue:
				i, ok := da.(int)
				if !ok || i != int(o.IntValue) {
					continue subscriptions
				}
			case *Filter_FloatValue:
				f, ok := da.(float64)
				if !ok || f != o.FloatValue {
					continue subscriptions
				}
			case *Filter_BoolValue:
				b, ok := da.(bool)
				if !ok || b != o.BoolValue {
					continue subscriptions
				}
			case *Filter_BytesValue:
				b, ok := da.([]byte)
				if !ok || len(b) != len(o.BytesValue) {
					continue subscriptions
				}
				for i := range b {
					z := b[i]
					x := o.BytesValue[i]
					if z != x {
						continue subscriptions
					}
				}
			case *Filter_TimeValue:
				s, ok := da.(time.Time)
				if !ok {
					continue subscriptions
				}
				if t := time.Unix(0, o.TimeValue); !t.Equal(s) {
					continue subscriptions
				}
			}
		}
		if d.Strict {
			var fm = make(map[string]struct{})
			if len(d.Filters) != len(args) {
				continue subscriptions
			}
			for _, filter := range d.Filters {
				if _, ok := fm[filter.GetKey()]; !ok {
					fm[filter.GetKey()] = struct{}{}
				}
			}
			for k := range args {
				if _, ok := fm[k]; !ok {
					continue subscriptions
				}
			}
		}
		parent := sub.
			Connection().
			Context().
			Value(router.RequestContext).(router.Context)
		child := parent.
			Child()
		params := graphql.Params{
			Schema:         v.schema,
			RequestString:  sub.Query(),
			VariableValues: sub.Variables(),
			OperationName:  sub.OperationName(),
			Context:        child.Context(),
			RootObject:     map[string]interface{}{"$subscription_payload$": d.Payload},
		}
		result := graphql.Do(params)
		if err := sub.Connection().Write(sub.ID(), &DataMessagePayload{
			Data:   result.Data,
			Errors: ErrorsFromGraphQLErrors(result.Errors),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (v *server) Subscriptions() Store {
	return v.tree
}

func (v *server) Handle(r router.Context, n *websocket.Conn) {
	u := fmt.Sprintf("%s-%s", r.MachineID(), r.ProcessID())
	c := &connection{
		machineID:     r.MachineID(),
		processID:     r.ProcessID(),
		identity:      r.Identity(),
		subscriptions: &subscriptions{subs: make(map[string]Subscription, 0)},
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

	if mt == websocket.TextMessage && message != nil && len(message) > 0 {

		raw := json.RawMessage{}
		msg := OperationMessage{Payload: &raw}

		if err := json.Unmarshal(message, &msg); err != nil {
			c.Close()
			return
		}

		switch msg.Type {

		// When the GraphQL WS connection is initiated, send an ACK back
		case gqlConnectionInit:

			data := new(InitMessagePayload)
			if err := json.Unmarshal(raw, &data); err != nil {
				c.Error(gqlUnknown, errors.New("Invalid GQL_CONNECTION_INIT payload"))
				return
			}
			if fn := v.hookauth; fn != nil {
				if err := fn(c, data.AuthToken); err != nil {
					c.Error(gqlUnknown, fmt.Errorf("Failed to authenticate user: %v", err))
					return
				}
			}
			c.Write(gqlUnknown, &OperationMessage{Type: gqlConnectionAck})

		// Handlers deal with starting operations
		case gqlStart:

			data := new(StartMessagePayload)
			if err := json.Unmarshal(raw, &data); err != nil {
				c.Error(msg.ID, errors.New("Invalid GQL_START payload"))
				return
			}

			if c.Subscriptions().Has(msg.ID) {
				c.Error(msg.ID, []error{errors.New("Cannot register subscription twice")})
				return
			}

			sub := &subscription{
				id:        msg.ID,
				query:     data.Query,
				opname:    data.OperationName,
				variables: data.Variables,
				args:      make(map[string]interface{}),
				ctx:       c,
			}

			if errs := validateSubscription(sub); len(errs) > 0 {
				c.Error(msg.ID, errs)
				return
			}

			doc, err := parser.Parse(parser.ParseParams{
				Source: sub.query,
			})
			if err != nil {
				c.Error(msg.ID, []error{err})
				return
			}

			if validation := graphql.ValidateDocument(&v.schema, doc, nil); !validation.IsValid {
				c.Error(msg.ID, ErrorsFromGraphQLErrors(validation.Errors))
				return
			}

			sub.doc = doc

			fields, args := getSubscriptionInfo(doc, data.Variables)
			sub.fields = fields
			sub.args = args

			if !c.Subscriptions().Add(sub, true) {
				c.Error(msg.ID, []error{errors.New("Unable to register subscription")})
				return
			}
			if !v.tree.Add(sub) {
				c.Error(msg.ID, []error{errors.New("Unable to register subscription")})
				return
			}

			c.Write(msg.ID, &OperationMessage{Type: gqlSubscriptionSuccess})

		// Handle all the stopping operations here
		case gqlStop:

			if id := strings.TrimSpace(msg.ID); len(id) == 0 {
				c.Error(gqlUnknown, errors.New("Invalid GQL_STOP payload"))
				return
			}

			if !c.Subscriptions().Has(msg.ID) {
				c.Error(msg.ID, []error{errors.New("Unable to find matching subscription")})
				return
			}

			sub := c.Subscriptions().Get(msg.ID)
			if !v.tree.Del(sub) || !c.Subscriptions().Del(sub) {
				c.Error(msg.ID, []error{errors.New("Unable to deregister subscription")})
				return
			}

		// When the GraphQL WS connection is terminated by the client,
		// close the connection
		case gqlConnectionTerminate:
			c.Close()

		// GraphQL WS protocol messages that are not handled
		default:
		}

	}
}

func (v *server) Authorize(f AuthorizeFunc) {
	v.hookauth = f
}

func (v *server) BeforeOpen(f HookFunc) {
	v.hookbo = f
}

func (v *server) AfterClose(f HookFunc) {
	v.hookac = f
}

func (v *server) String() string {
	return "Hyper::GraphQLSubscription"
}
