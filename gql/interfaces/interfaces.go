package interfaces

import (
	"context"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/router"
)

// ResolveHandler for field resolve
type ResolveHandler func(Resolver) (interface{}, error)

// Schema for GraphQL
type Schema interface {
	Query(Object) Schema
	Mutation(Object) Schema
	Subscription(Object) Schema
	Config() SchemaConfig
}

// SchemaConfig interface
type SchemaConfig interface {
	Query() Object
	Mutation() Object
	Subscription() Object
	Schema() graphql.Schema
}

// Object for GraphQL
type Object interface {
	Description(string) Object
	Fields(...Field) Object
	RecursiveFields(...Field) Object
	Args(...Argument) Object
	Config() ObjectConfig
}

// ObjectConfig interface
type ObjectConfig interface {
	Name() string
	Description() string
	Fields() []Field
	RecursiveFields() []Field
	Args() []Argument
	Output() *graphql.Object
	HasOutput() bool
	Input() *graphql.InputObject
	HasInput() bool
}

// Field for GraphQL
type Field interface {
	Description(string) Field
	DeprecationReason(string) Field
	Type(interface{}) Field
	Args(...Argument) Field
	Resolve(ResolveHandler) Field
	Config() FieldConfig
}

// FieldConfig interface
type FieldConfig interface {
	Name() string
	Description() string
	DeprecationReason() string
	Type() graphql.Output
	Args() []Argument
	Field() *graphql.Field
}

// Argument for GraphQL
type Argument interface {
	Description(string) Argument
	Type(interface{}) Argument
	Default([]byte) Argument
	Require(bool) Argument
	Config() ArgumentConfig
}

// ArgumentConfig interface
type ArgumentConfig interface {
	Name() string
	Description() string
	Type() graphql.Input
	Object() Object
	Default() []byte
	Require() bool
	ArgumentConfig() *graphql.ArgumentConfig
	InputObjectFieldConfig() *graphql.InputObjectFieldConfig
}

// Resolver interface
type Resolver interface {
	Context() Context
	Params() graphql.ResolveParams
	Source() interface{}
	Arg(string) (Value, error)
	MustArg(string) Value
}

// Value for GraphQL
type Value interface {
	In(string) Value
	Key() string
	Val() []byte
	Has() bool
	MustInt() int
	MustI32() int32
	MustI64() int64
	MustU32() uint32
	MustU64() uint64
	MustF32() float32
	MustF64() float64
	MustBool() bool
	MustTime() time.Time
	MustArray() []interface{}
	Any() interface{}
	String() string
}

// Context interface
type Context interface {
	Deadline() (time.Time, bool)
	Done() <-chan struct{}
	Err() error
	Value(interface{}) interface{}
	Identity() router.Identity
	MachineID() string
	ProcessID() string
	Context() context.Context
	Req() *http.Request
	Res() http.ResponseWriter
	Client() router.Client
	Cache() router.CacheAdaptor
	Message() router.MessageAdaptor
	DataLoader(interface{}) router.DataLoaderAdaptor
	KV() router.KV
	Cookie() router.Cookie
	Header() router.Header
	MustParam(s string) router.Value
	MustQuery(s string) router.Value
	MustBody(s string) router.Value
	Param(s string) (router.Value, error)
	Query(s string) (router.Value, error)
	Body(s string) (router.Value, error)
	File(s string) []byte
	Abort()
	IsAborted() bool
	HasErrors() bool
	Errors() []error
	GraphQLError(error)
	Native() router.Context
}
