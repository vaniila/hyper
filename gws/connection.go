package gws

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vaniila/hyper/router"
)

type connection struct {
	machineID, processID string
	ctx                  context.Context
	identity             Identity
	subscriptions        Subscriptions
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

func (v *connection) Connection() *websocket.Conn {
	return v.conn
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

func (v *connection) Write(id string, o interface{}) error {
	var msg *OperationMessage
	switch d := o.(type) {
	case *OperationMessage:
		msg = d
	default:
		msg = &OperationMessage{
			Type:    gqlData,
			Payload: o,
		}
	}
	msg.ID = id
	return v.conn.WriteMessage(websocket.TextMessage, msg.Marshal())
}

func (v *connection) Error(id string, errs interface{}) error {
	var dat interface{}
	switch errs.(type) {
	case []error, error:
		dat = errs
	default:
		return nil
	}
	o := &OperationMessage{
		ID:      id,
		Type:    gqlError,
		Payload: dat,
	}
	return v.conn.WriteMessage(websocket.TextMessage, o.Marshal())
}

func (v *connection) Close() error {
	return v.conn.Close()
}

func (v *connection) BeforeOpen() {
}

func (v *connection) AfterClose() {
	// clean up subscriptions
	for _, c := range v.Subscriptions().List() {
		v.server.Subscriptions().Del(c)
		v.Subscriptions().Del(c)
	}
}
