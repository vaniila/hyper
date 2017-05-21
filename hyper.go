package hyper

import (
	"github.com/samuelngs/hyper/engine"
	"github.com/samuelngs/hyper/websocket"
)

// New creates a hyper server
func New(opts ...Option) *Hyper {
	o := newOptions(opts...)
	e := engine.New(
		engine.ID(o.ID),
		engine.Addr(o.Addr),
		engine.Proto(o.Protocol),
		engine.Cache(o.Cache),
		engine.Message(o.Message),
		engine.Router(o.Router),
	)
	w := websocket.New(
		websocket.ID(o.ID),
		websocket.Cache(o.Cache),
		websocket.Message(o.Message),
		websocket.Router(o.Router),
	)
	return &Hyper{
		id:        o.ID,
		addr:      o.Addr,
		cache:     o.Cache,
		router:    o.Router,
		message:   o.Message,
		engine:    e,
		websocket: w,
	}
}
