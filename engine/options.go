package engine

import (
	"crypto/rand"
	"fmt"

	"github.com/vaniila/hyper/cache"
	"github.com/vaniila/hyper/dataloader"
	"github.com/vaniila/hyper/gws"
	"github.com/vaniila/hyper/logger"
	"github.com/vaniila/hyper/message"
	"github.com/vaniila/hyper/router"
	"github.com/vaniila/hyper/websocket"
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

	// Logger
	Logger logger.Service

	// GraphQL subscription server
	GQLSubscription gws.Service

	// DataLoader
	DataLoader dataloader.Service

	// Router
	Router router.Service

	// Websocket service
	Websocket websocket.Service

	// TraceID customize function
	TraceID func() string

	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penality.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins []string

	// AllowOriginFunc is a custom function to validate the origin. It take the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	AllowOriginFunc func(origin string) bool

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (GET and POST)
	AllowedMethods []string

	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge int

	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool
}

func newID() string {
	b := new([16]byte)
	rand.Read(b[:])
	b[8] = (b[8] | 0x40) & 0x7F
	b[6] = (b[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func newOptions(opts ...Option) Options {
	opt := Options{
		ID:       newID(),
		Addr:     ":0",
		Protocol: HTTP,
		TraceID:  newID,
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
	if opt.Message == nil {
		opt.DataLoader = dataloader.New(
			dataloader.ID(opt.ID),
		)
	}
	if opt.TraceID == nil {
		opt.TraceID = newID
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

// Logger to set custom logger
func Logger(l logger.Service) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// GQLSubscription to bind graphql subscription interface to engine server
func GQLSubscription(s gws.Service) Option {
	return func(o *Options) {
		o.GQLSubscription = s
	}
}

// DataLoader server
func DataLoader(d dataloader.Service) Option {
	return func(o *Options) {
		o.DataLoader = d
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

// TraceID to set trace id generator function
func TraceID(f func() string) Option {
	return func(o *Options) {
		o.TraceID = f
	}
}

// AllowedOrigins to set allowed origins
func AllowedOrigins(a []string) Option {
	return func(o *Options) {
		o.AllowedOrigins = append(o.AllowedOrigins, a...)
	}
}

// AllowOriginFunc to add custom origin handler
func AllowOriginFunc(f func(string) bool) Option {
	return func(o *Options) {
		o.AllowOriginFunc = f
	}
}

// AllowedMethods to set allowed methods
func AllowedMethods(a []string) Option {
	return func(o *Options) {
		o.AllowedMethods = append(o.AllowedMethods, a...)
	}
}

// AllowedHeaders to set allowed headers
func AllowedHeaders(a []string) Option {
	return func(o *Options) {
		o.AllowedHeaders = append(o.AllowedMethods, a...)
	}
}

// ExposedHeaders to set exposed headers
func ExposedHeaders(a []string) Option {
	return func(o *Options) {
		o.ExposedHeaders = append(o.AllowedMethods, a...)
	}
}

// AllowCredentials to set allow credentials header
func AllowCredentials(b bool) Option {
	return func(o *Options) {
		o.AllowCredentials = b
	}
}

// MaxAge to set max age value
func MaxAge(i int) Option {
	return func(o *Options) {
		o.MaxAge = i
	}
}

// OptionsPassthrough to set options pass through value
func OptionsPassthrough(b bool) Option {
	return func(o *Options) {
		o.OptionsPassthrough = b
	}
}
