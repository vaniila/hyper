package router

import (
	"context"
	"net/http"

	"github.com/samuelngs/hyper/router/cookie"
	"github.com/samuelngs/hyper/tracer"
)

// HandlerFunc type
type HandlerFunc func(Context)

// HandlerFuncs type
type HandlerFuncs []HandlerFunc

// Context interface
type Context interface {
	MachineID() string
	ProcessID() string
	Context() context.Context
	Req() *http.Request
	Res() http.ResponseWriter
	Client() Client
	Cache() Cache
	Cookie() Cookie
	Header() Header
	MustParam(s string) Value
	MustQuery(s string) Value
	MustBody(s string) Value
	Param(s string) (Value, error)
	Query(s string) (Value, error)
	Body(s string) (Value, error)
	File(s string) []byte
	Tracer() tracer.Tracer
	Recover() error
	Abort()
	IsAborted() bool
	Write(b []byte) Context
	Error(error) Context
	Json(o interface{}) Context
	Status(code int) Context
}

// Cache interface
type Cache interface {
	Set(key string, data []byte)
	Get(key string) []byte
}

// Value interface
type Value interface {
	Key() string
	Val() []byte
	String() string
}

// Header interface
type Header interface {
	Set(key string, val string)
	Get(key string) (Value, error)
	MustGet(key string) Value
}

// Cookie interface
type Cookie interface {
	Set(key string, val string, opts ...cookie.Option)
	Get(key string) (Value, error)
	MustGet(key string) Value
}

// Client interface
type Client interface {
	IP() string
	Host() string
	Device() Device
	Protocol() string
}

// Device interface
type Device interface {
	Useragent() Useragent
	Hardware() Hardware
	OS() OS
}

// Useragent interface
type Useragent interface {
	Family() string
	Major() string
	Minor() string
	Patch() string
	String() string
}

// OS interface
type OS interface {
	Family() string
	Major() string
	Minor() string
	Patch() string
}

// Hardware interface
type Hardware interface {
	Family() string
	Brand() string
	Model() string
}
