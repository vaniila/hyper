package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/vaniila/hyper/cache"
	"github.com/vaniila/hyper/gws"
	"github.com/vaniila/hyper/logger"
	"github.com/vaniila/hyper/message"
	"github.com/vaniila/hyper/router"
	"github.com/vaniila/hyper/sync"
)

type server struct {
	id       string
	sync     sync.Service
	gws      gws.Service
	cache    cache.Service
	message  message.Service
	logger   logger.Service
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
	switch conn.Subprotocol() {
	case "graphql-ws":
		if v.gws != nil {
			v.gws.Handle(c, conn)
		}
	default:
		if v.sync != nil {
			v.sync.Handle(c, conn)
		}
	}
}

func (v *server) String() string {
	return "Hyper::Websocket"
}
