package engine

import (
	"crypto/rand"
	"fmt"

	"github.com/samuelngs/hyper/cache"
	"github.com/samuelngs/hyper/message"
	"github.com/samuelngs/hyper/router"
	"github.com/samuelngs/hyper/websocket"
)

// Option func
type Option func(*Options)

// Options is the engine server options
type Options struct {

	// engine server unique id
	ID string

	// server bind to address [host:port]
	Addr string

	// HTTP protocol 1.1 / 2.0
	Protocol Protocol

	// Cache
	Cache cache.Service

	// Message broker
	Message message.Service

	// Router
	Router router.Service

	// Websocket service
	Websocket websocket.Service
}

func newID() string {
	b := new([16]byte)
	rand.Read(b[:])
	b[8] = (b[8] | 0x40) & 0x7F
	b[6] = (b[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func newOptions(opts ...Option) Options {
	opt := Options{
		ID:       newID(),
		Addr:     ":0",
		Protocol: HTTP,
	}
	for _, o := range opts {
		o(&opt)
	}
	if opt.Cache == nil {
		opt.Cache = cache.New(
			cache.ID(opt.ID),
		)
	}
	if opt.Router == nil {
		opt.Router = router.New(
			router.ID(opt.ID),
		)
	}
	if opt.Message == nil {
		opt.Message = message.New(
			message.ID(opt.ID),
		)
	}
	return opt
}

// ID to change server reference id
func ID(s string) Option {
	return func(o *Options) {
		o.ID = s
	}
}

// Addr to change server bind address
func Addr(s string) Option {
	return func(o *Options) {
		o.Addr = s
	}
}

// Proto to change http network protocol
func Proto(p Protocol) Option {
	return func(o *Options) {
		o.Protocol = p
	}
}

// Cache to bind cache interface to engine server
func Cache(c cache.Service) Option {
	return func(o *Options) {
		o.Cache = c
	}
}

// Message broker
func Message(m message.Service) Option {
	return func(o *Options) {
		o.Message = m
	}
}

// Router to bind router to engine server
func Router(r router.Service) Option {
	return func(o *Options) {
		o.Router = r
	}
}

// Websocket to bind websocket server
func Websocket(w websocket.Service) Option {
	return func(o *Options) {
		o.Websocket = w
	}
}
