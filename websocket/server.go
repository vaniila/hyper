package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/samuelngs/hyper/cache"
	"github.com/samuelngs/hyper/message"
	"github.com/samuelngs/hyper/router"
	"github.com/samuelngs/hyper/sync"
)

type server struct {
	id       string
	sync     sync.Service
	cache    cache.Service
	message  message.Service
	upgrader websocket.Upgrader
}

func (v *server) Start() error {
	return nil
}

func (v *server) Stop() error {
	return nil
}

func (v *server) Handle(c router.Context) {
	conn, err := v.upgrader.Upgrade(c.Res(), c.Req(), nil)
	if err != nil {
		return
	}
	defer conn.Close()
	if v.sync != nil {
		v.sync.Handle(c, conn)
	}
}

func (v *server) String() string {
	return "Hyper::Websocket"
}
