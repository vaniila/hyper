package gql

import (
	"context"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/opentracing/opentracing-go"
	"github.com/vaniila/hyper/router"
)

// ObjectInitializer func
type ObjectInitializer func(Object)

// FieldInitializer func
type FieldInitializer func(Field)

// ArgumentInitializer func
type ArgumentInitializer func(Argument)

// ResolveHandler for field resolve
type ResolveHandler func(Resolver) (interface{}, error)

// ScalarSerializeHandler for scalar serialize resolver
type ScalarSerializeHandler func(interface{}) (interface{}, error)

// ScalarParseValueHandler for scalar parse value resolver
type ScalarParseValueHandler func(interface{}) (interface{}, error)

// ScalarParseLiteralHandler for scalar parse literal resolver
type ScalarParseLiteralHandler func(ast.Value) (interface{}, error)

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

// Union for GraphQL
type Union interface {
	Description(string) Union
	Resolve(interface{}, Object) Union
	Config() UnionConfig
}

// UnionConfig interface
type UnionConfig interface {
	Name() string
	Description() string
	Union() *graphql.Union
}

// Scalar for GraphQL
type Scalar interface {
	Description(string) Scalar
	Serialize(ScalarSerializeHandler) Scalar
	ParseValue(ScalarParseValueHandler) Scalar
	ParseLiteral(ScalarParseLiteralHandler) Scalar
	Config() ScalarConfig
}

// ScalarConfig interface
type ScalarConfig interface {
	Name() string
	Description() string
	Serialize(interface{}) (interface{}, error)
	ParseValue(interface{}) (interface{}, error)
	ParseLiteral(ast.Value) (interface{}, error)
}

// Enum for GraphQL
type Enum interface {
	Description(string) Enum
	Values(...EnumValue) Enum
	Config() EnumConfig
}

// EnumConfig interface
type EnumConfig interface {
	Name() string
	Description() string
	Values() []EnumValue
	Enum() *graphql.Enum
}

// EnumValue for GraphQL
type EnumValue interface {
	Description(string) EnumValue
	Is(interface{}) EnumValue
	Deprecation(string) EnumValue
	Config() EnumValueConfig
}

// EnumValueConfig interface
type EnumValueConfig interface {
	Name() string
	Description() string
	Is() interface{}
	Deprecation() string
}

// Object for GraphQL
type Object interface {
	Description(string) Object
	Fields(...Field) Object
	Args(...Argument) Object
	Init(ObjectInitializer) Object
	Config() ObjectConfig
}

// ObjectConfig interface
type ObjectConfig interface {
	Name() string
	Description() string
	Fields() []Field
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
	Init(FieldInitializer) Field
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
	Default(interface{}) Argument
	Require(bool) Argument
	Init(ArgumentInitializer) Argument
	Config() ArgumentConfig
}

// ArgumentConfig interface
type ArgumentConfig interface {
	Name() string
	Description() string
	Type() graphql.Input
	Object() Object
	Default() interface{}
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
	Logger() router.Logger
	GQLSubscription() router.GQLSubscriptionAdaptor
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
	StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span
	Tracer() opentracing.Tracer
	Abort()
	IsAborted() bool
	HasErrors() bool
	Errors() []error
	GraphQLError(error)
	Native() router.Context
}
