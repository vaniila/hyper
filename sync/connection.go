package sync

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/samuelngs/hyper/router"
)

type connection struct {
	machineID, processID string
	ctx                  context.Context
	identity             *identity
	subscriptions        *subscriptions
	req                  *http.Request
	res                  http.ResponseWriter
	client               router.Client
	cookie               router.Cookie
	header               router.Header
	cache                CacheAdaptor
	message              MessageAdaptor
	server               Service
	conn                 *websocket.Conn
}

func (v *connection) MachineID() string {
	return v.machineID
}

func (v *connection) ProcessID() string {
	return v.processID
}

func (v *connection) Identity() Identity {
	return v.identity
}

func (v *connection) Subscriptions() Subscriptions {
	return v.subscriptions
}

func (v *connection) Context() context.Context {
	return v.ctx
}

func (v *connection) Req() *http.Request {
	return v.req
}

func (v *connection) Res() http.ResponseWriter {
	return v.res
}

func (v *connection) Client() router.Client {
	return v.client
}

func (v *connection) Cookie() router.Cookie {
	return v.cookie
}

func (v *connection) Header() router.Header {
	return v.header
}

func (v *connection) Cache() CacheAdaptor {
	return v.cache
}

func (v *connection) Message() MessageAdaptor {
	return v.message
}

func (v *connection) Write(p *Packet) error {
	b, err := proto.Marshal(p)
	if err != nil {
		return err
	}
	return v.conn.WriteMessage(websocket.BinaryMessage, b)
}

func (v *connection) Close() error {
	return v.conn.Close()
}

func (v *connection) BeforeOpen() {
}

func (v *connection) AfterClose() {
	// clean up channels
	// close channels that has no subscribers in it
	for _, c := range v.Subscriptions().List() {
		c.Unsubscribe(v)
		l := len(c.NodeSubscribers())
		if l == 0 {
			c.Namespace().Channels().Del(c.Name())
		}
	}
}
