package sync

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/samuelngs/hyper/message"
	"github.com/samuelngs/hyper/router"
)

// HookFunc
type HookFunc func(Context)

// AuthorizeFunc
type AuthorizeFunc func(string, Context) error

// HandlerFunc type
type HandlerFunc func([]byte, Channel, Context)

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
	Publish(*Distribution) error
	Subscribe(*Distribution) error
	Handle(router.Context, *websocket.Conn)
	BeforeOpen(HookFunc)
	AfterClose(HookFunc)
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
	Channels() Channels
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
	Channels() Channels
	Doc() string
}

type Channels interface {
	Namespace() Namespace
	Has(string) bool
	Get(string) Channel
	Add(string) Channels
	Del(string) Channels
	List() map[string]Channel
	Len() int
}

// Channel interface
type Channel interface {
	Namespace() Namespace
	Name() string
	NodeSubscribers() []Context
	Has(Context) bool
	Subscribe(Context) Channel
	Unsubscribe(Context) Channel
	Write(*Packet, ...*Condition) error
	BeforeOpen()
	AfterClose()
}

// Context interface
type Context interface {
	Identity() Identity
	MachineID() string
	ProcessID() string
	Context() context.Context
	Req() *http.Request
	Res() http.ResponseWriter
	Client() router.Client
	Cookie() router.Cookie
	Header() router.Header
	Subscriptions() Subscriptions
	Cache() CacheAdaptor
	Message() MessageAdaptor
	Write(*Packet) error
	Close() error
	BeforeOpen()
	AfterClose()
}

// Identity interface
type Identity interface {
	HasID() bool
	GetID() int
	SetID(int)
	HasKey() bool
	GetKey() string
	SetKey(string)
}

// Subscriptions interface
type Subscriptions interface {
	Has(string, string) bool
	Add(Channel) Subscriptions
	Del(Channel) Subscriptions
	List() []Channel
}

// Cache interface
type CacheAdaptor interface {
	Set(key []byte, data []byte, ttl time.Duration) error
	Get(key []byte) ([]byte, error)
}

// Message broker interface
type MessageAdaptor interface {
	Emit([]byte, []byte) error
	Listen([]byte, message.Handler) message.Close
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
		nsmap:      make(map[string]Namespace),
		conns:      make(map[string]Context),
	}
	return s
}
