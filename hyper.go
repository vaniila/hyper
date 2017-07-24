package hyper

import (
	"github.com/vaniila/hyper/engine"
	"github.com/vaniila/hyper/websocket"
)

// New creates a hyper server
func New(opts ...Option) *Hyper {
	o := newOptions(opts...)
	w := websocket.New(
		websocket.ID(o.ID),
		websocket.Sync(o.Sync),
		websocket.Cache(o.Cache),
		websocket.Message(o.Message),
		websocket.Router(o.Router),
		websocket.EnableCompression(true),
	)
	e := engine.New(
		engine.ID(o.ID),
		engine.Addr(o.Addr),
		engine.Proto(o.Protocol),
		engine.Cache(o.Cache),
		engine.Message(o.Message),
		engine.Router(o.Router),
		engine.Websocket(w),
		engine.AllowedOrigins(o.AllowedOrigins),
		engine.AllowOriginFunc(o.AllowOriginFunc),
		engine.AllowedMethods(o.AllowedMethods),
		engine.AllowedHeaders(o.AllowedHeaders),
		engine.ExposedHeaders(o.ExposedHeaders),
		engine.AllowCredentials(o.AllowCredentials),
		engine.MaxAge(o.MaxAge),
		engine.OptionsPassthrough(o.OptionsPassthrough),
	)
	return &Hyper{
		id:        o.ID,
		addr:      o.Addr,
		cache:     o.Cache,
		message:   o.Message,
		sync:      o.Sync,
		router:    o.Router,
		engine:    e,
		websocket: w,
	}
}
