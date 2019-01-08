package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/vaniila/hyper/router"
)

// Service interface
type Service interface {
	Start() error
	Stop() error
	Handle(router.Context)
	String() string
}

// New creates engine server
func New(opts ...Option) Service {
	o := newOptions(opts...)
	s := &server{
		id:      o.ID,
		sync:    o.Sync,
		gws:     o.GQLSubscription,
		cache:   o.Cache,
		message: o.Message,
		upgrader: websocket.Upgrader{
			Subprotocols:      []string{"graphql-ws", "hyper-ws"},
			HandshakeTimeout:  o.HandshakeTimeout,
			ReadBufferSize:    o.ReadBufferSize,
			WriteBufferSize:   o.WriteBufferSize,
			CheckOrigin:       o.CheckOrigin,
			EnableCompression: o.EnableCompression,
		},
	}
	return s
}
