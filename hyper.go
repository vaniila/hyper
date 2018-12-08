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
		websocket.GQLSubscription(o.GQLSubscription),
		websocket.Cache(o.Cache),
		websocket.Message(o.Message),
		websocket.Router(o.Router),
		websocket.Logger(o.Logger),
		websocket.EnableCompression(true),
	)
	e := engine.New(
		engine.ID(o.ID),
		engine.Addr(o.Addr),
		engine.Proto(o.Protocol),
		engine.Cache(o.Cache),
		engine.Message(o.Message),
		engine.Logger(o.Logger),
		engine.GQLSubscription(o.GQLSubscription),
		engine.DataLoader(o.DataLoader),
		engine.Router(o.Router),
		engine.Websocket(w),
		engine.TraceID(o.TraceID),
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
		id:         o.ID,
		addr:       o.Addr,
		cache:      o.Cache,
		message:    o.Message,
		logger:     o.Logger,
		dataloader: o.DataLoader,
		sync:       o.Sync,
		gws:        o.GQLSubscription,
		router:     o.Router,
		engine:     e,
		websocket:  w,
	}
}
