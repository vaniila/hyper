package websocket

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/samuelngs/hyper/cache"
	"github.com/samuelngs/hyper/message"
	"github.com/samuelngs/hyper/router"
	"github.com/samuelngs/hyper/sync"
)

// Option func
type Option func(*Options)

// Options is the engine server options
type Options struct {

	// engine server unique id
	ID string

	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration

	// ReadBufferSize and WriteBufferSize specify I/O buffer sizes. If a buffer
	// size is zero, then a default value of 4096 is used. The I/O buffer sizes
	// do not limit the size of the messages that can be sent or received.
	ReadBufferSize, WriteBufferSize int

	// CheckOrigin returns true if the request Origin header is acceptable. If
	// CheckOrigin is nil, the host in the Origin header must not be set or
	// must match the host of the request.
	CheckOrigin func(r *http.Request) bool

	// EnableCompression specify if the server should attempt to negotiate per
	// message compression (RFC 7692). Setting this value to true does not
	// guarantee that compression will be supported. Currently only "no context
	// takeover" modes are supported.
	EnableCompression bool

	// Sync engine server
	Sync sync.Service

	// Cache server
	Cache cache.Service

	// Message broker
	Message message.Service

	// Router
	Router router.Service
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
		ID: newID(),
	}
	for _, o := range opts {
		o(&opt)
	}
	if opt.CheckOrigin == nil {
		opt.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}
	return opt
}

// ID to change server reference id
func ID(s string) Option {
	return func(o *Options) {
		o.ID = s
	}
}

// HandshakeTimeout to set server handshake timeout
func HandshakeTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.HandshakeTimeout = d
	}
}

// ReadBufferSize to set server read buffer size
func ReadBufferSize(i int) Option {
	return func(o *Options) {
		o.ReadBufferSize = i
	}
}

// WriteBufferSize to set server write buffer size
func WriteBufferSize(i int) Option {
	return func(o *Options) {
		o.WriteBufferSize = i
	}
}

// CheckOrigin to set origin checking function
func CheckOrigin(f func(r *http.Request) bool) Option {
	return func(o *Options) {
		o.CheckOrigin = f
	}
}

// EnableCompression to set compression settings for websocket
func EnableCompression(b bool) Option {
	return func(o *Options) {
		o.EnableCompression = b
	}
}

// Sync to bind sync interface to websocket server
func Sync(s sync.Service) Option {
	return func(o *Options) {
		o.Sync = s
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
