package router

import (
	"context"
	"net/http"

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
	Param() Value
	Query() Value
	Body() Value
	File(s string) []byte
	Tracer() tracer.Tracer
	Abort()
	IsAborted() bool
	Write(b []byte) Context
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
	String() string
}

// Header interface
type Header interface {
	Set(key string, val string)
	Get(key string) Value
}

// Cookie interface
type Cookie interface {
	Set(key string, val string)
	Get(key string) Value
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
