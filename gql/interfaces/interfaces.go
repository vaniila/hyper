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
	Compile() graphql.Schema
}

// Object for GraphQL
type Object interface {
	Name(string) Object
	Description(string) Object
	Fields(...Field) Object
	Args(...Argument) Object
	ToObject() *graphql.Object
	ToInputObject() *graphql.InputObject
	ExportFields() []Field
	ExportArgs() []Argument
}

// Field for GraphQL
type Field interface {
	Name(string) Field
	Description(string) Field
	DeprecationReason(string) Field
	Type(interface{}) Field
	Args(...Argument) Field
	Resolve(ResolveHandler) Field
	Compile() *graphql.Field
}

// Argument for GraphQL
type Argument interface {
	Name(string) Argument
	Description(string) Argument
	Type(interface{}) Argument
	Default([]byte) Argument
	Require(bool) Argument
	InputObject() Object
	ToArgumentConfig() (string, *graphql.ArgumentConfig)
	ToInputObjectFieldConfig() (string, *graphql.InputObjectFieldConfig)
}

// Resolver interface
type Resolver interface {
	Context() Context
	Params() graphql.ResolveParams
	Source() interface{}
	Arg(string) (Value, error)
	MustArg(string) Value
}

// Type for GraphQL
type Type interface {
	Name(string) Type
	Description(string) Type
	Fields(...Field) Type
	Input() graphql.Input
	Output() graphql.Output
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
}
