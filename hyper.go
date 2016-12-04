package hyper

import (
	"github.com/samuelngs/hyper/cache"
	"github.com/samuelngs/hyper/engine"
	"github.com/samuelngs/hyper/router"
	"github.com/samuelngs/hyper/websocket"
)

// New creates a hyper server
func New(opts ...Option) *Hyper {
	o := newOptions(opts...)
	c := cache.New(
		cache.ID(o.ID),
	)
	r := router.New(
		router.ID(o.ID),
	)
	e := engine.New(
		engine.ID(o.ID),
		engine.Addr(o.Addr),
		engine.Proto(o.Protocol),
		engine.Router(r),
	)
	w := websocket.New(
		websocket.ID(o.ID),
		websocket.Router(r),
	)
	return &Hyper{
		id:        o.ID,
		addr:      o.Addr,
		cache:     c,
		router:    r,
		engine:    e,
		websocket: w,
	}
}
