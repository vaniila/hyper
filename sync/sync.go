package sync

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/samuelngs/hyper/router"
)

// AuthorizeFunc
type AuthorizeFunc func(Context) error

// HandlerFunc type
type HandlerFunc func([]byte, Context)

// HandlerFuncs type
type HandlerFuncs []HandlerFunc

// Handler interface
type Handler interface {
	Name() string
	Func() HandlerFunc
}

// Service interface
type Service interface {
	Start() error
	Stop() error
	Namespace(string) Namespace
	Namespaces() []Namespace
	Publish(Packet) error
	Subscribe(Packet) error
	Handle(router.Context, *websocket.Conn)
	String() string
}

// Namespace interface
type Namespace interface {
	Alias(...string) Namespace
	Name(string) Namespace
	Summary(string) Namespace
	Doc(string) Namespace
	Authorize(AuthorizeFunc) Namespace
	Middleware(...HandlerFunc) Namespace
	Handle(string, HandlerFunc) Namespace
	Catch(HandlerFunc) Namespace
	Config() NamespaceConfig
}

type NamespaceConfig interface {
	Namespace() string
	Aliases() []string
	Name() string
	Summary() string
	Authorize() AuthorizeFunc
	Middlewares() []HandlerFunc
	Handlers() []Handler
	Catch() HandlerFunc
	Doc() string
}

// Packet interface
type Packet interface {
	Direction() int
	Type() int
	Namespace() string
	Channel() string
	Message() []byte
	Bytes() ([]byte, error)
}

// Channel interface
type Channel interface {
	Namespace() Namespace
	Name() string
	Subscribed() bool
	BeforeOpen()
	AfterClose()
}

// Context interface
type Context interface {
	MachineID() string
	ProcessID() string
	Channels() []Channel
	Context() context.Context
	Req() *http.Request
	Res() http.ResponseWriter
	Client() router.Client
	Cookie() router.Cookie
	Header() router.Header
	Cache() CacheAdaptor
	Message() MessageAdaptor
	Write(Packet) error
	Close() error
	BeforeOpen()
	AfterClose()
}

// Cache interface
type CacheAdaptor interface {
	Set(key []byte, data []byte, ttl time.Duration) error
	Get(key []byte) ([]byte, error)
}

// Message broker interface
type MessageAdaptor interface {
	Emit([]byte, []byte) error
	Listen([]byte) (<-chan []byte, chan<- struct{}, error)
}

// New creates engine server
func New(opts ...Option) Service {
	o := newOptions(opts...)
	s := &server{
		id:         o.ID,
		topic:      o.Topic,
		cache:      o.Cache,
		message:    o.Message,
		namespaces: make([]Namespace, 0),
		conns:      make(map[string]Context),
	}
	return s
}
