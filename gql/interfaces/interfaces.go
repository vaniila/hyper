package interfaces

import (
	"context"
	"net/http"

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
	Compile() *graphql.Object
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
	Type(graphql.Input) Argument
	Default([]byte) Argument
	Require(bool) Argument
	Compile() (string, *graphql.ArgumentConfig)
}

// Resolver interface
type Resolver interface {
	Context() Context
	Params() graphql.ResolveParams
	Source() interface{}
	Arg(string) (router.Value, error)
	MustArg(string) router.Value
}

// Type for GraphQL
type Type interface {
	Name(string) Type
	Description(string) Type
	Fields(...Field) Type
	Input() graphql.Input
	Output() graphql.Output
}

// Context interface
type Context interface {
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
