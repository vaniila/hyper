package websocket

import (
	"github.com/samuelngs/hyper/cache"
	"github.com/samuelngs/hyper/message"
)

type server struct {
	id      string
	cache   cache.Service
	message message.Service
}

func (v *server) Start() error {
	return nil
}

func (v *server) Stop() error {
	return nil
}

func (v *server) String() string {
	return "Hyper::Websocket"
}
