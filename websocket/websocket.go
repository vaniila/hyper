package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/samuelngs/hyper/router"
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
		cache:   o.Cache,
		message: o.Message,
		upgrader: websocket.Upgrader{
			HandshakeTimeout:  o.HandshakeTimeout,
			ReadBufferSize:    o.ReadBufferSize,
			WriteBufferSize:   o.WriteBufferSize,
			CheckOrigin:       o.CheckOrigin,
			EnableCompression: o.EnableCompression,
		},
	}
	return s
}
