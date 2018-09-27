package gws

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/vaniila/hyper/logger"
	"github.com/vaniila/hyper/message"
	"github.com/vaniila/hyper/router"
)

// HookFunc type
type HookFunc func(Context)

// AuthorizeFunc type
type AuthorizeFunc func(Context, string) error

// HandlerFunc type
type HandlerFunc func([]byte, Context)

// Service interface
type Service interface {
	Start() error
	Stop() error
	Publish(*Distribution) error
	Subscribe(*Distribution) error
	Subscriptions() Store
	Handle(router.Context, *websocket.Conn)
	Schema(graphql.Schema)
	Adaptor() router.GQLSubscriptionAdaptor
	Authorize(AuthorizeFunc)
	BeforeOpen(HookFunc)
	AfterClose(HookFunc)
	String() string
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
	Connection() *websocket.Conn
	Cache() CacheAdaptor
	Message() MessageAdaptor
	Logger() LoggerAdaptor
	Write(string, interface{}) error
	Error(string, interface{}) error
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

// Store interface
type Store interface {
	Has(Subscription) bool
	Add(Subscription, ...bool) bool
	Del(Subscription) bool
	Get(...string) []Subscription
}

// Subscriptions interface
type Subscriptions interface {
	Has(string) bool
	Add(Subscription, ...bool) bool
	Del(Subscription) bool
	Get(string) Subscription
	List() []Subscription
	Len() int
}

// Subscription interface
type Subscription interface {
	ID() string
	Query() string
	Variables() map[string]interface{}
	Arguments() map[string]interface{}
	OperationName() string
	Document() *ast.Document
	Fields() []string
	Connection() Context
}

// CacheAdaptor interface
type CacheAdaptor interface {
	Set(key []byte, data []byte, ttl time.Duration) error
	Get(key []byte) ([]byte, error)
}

// MessageAdaptor broker interface
type MessageAdaptor interface {
	Emit([]byte, []byte) error
	Listen([]byte, message.Handler) message.Close
}

// LoggerAdaptor logging interface
type LoggerAdaptor interface {
	Debug(string, ...logger.Field)
	Info(string, ...logger.Field)
	Warn(string, ...logger.Field)
	Error(string, ...logger.Field)
	Fatal(string, ...logger.Field)
	Panic(string, ...logger.Field)
}

// New creates engine server
func New(opts ...Option) Service {
	o := newOptions(opts...)
	s := &server{
		id:      o.ID,
		topic:   o.Topic,
		cache:   o.Cache,
		message: o.Message,
		logger:  o.Logger,
		schema:  o.Schema,
		conns:   make(map[string]Context),
		tree:    &tree{state: make(map[string][]Subscription)},
	}
	s.adaptor = &adaptor{s}
	return s
}
