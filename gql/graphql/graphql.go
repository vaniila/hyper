package graphql

import (
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql"
	"github.com/vaniila/hyper/gql/argument"
	"github.com/vaniila/hyper/gql/enum"
	"github.com/vaniila/hyper/gql/field"
	"github.com/vaniila/hyper/gql/object"
	"github.com/vaniila/hyper/gql/schema"
	"github.com/vaniila/hyper/gql/union"
)

// builtin graphql scalars
var (
	Int      = graphql.Int
	Float    = graphql.Float
	String   = graphql.String
	Boolean  = graphql.Boolean
	ID       = graphql.ID
	DateTime = graphql.DateTime
)

type (
	// Resolver alias
	Resolver = gql.Resolver
	// Context alias
	Context = gql.Context
)

// Schema creates new schema
func Schema(opts ...schema.Option) graphql.Schema {
	return schema.New(opts...).Config().Schema()
}

// Query option
func Query(c gql.Object) schema.Option {
	return schema.Query(c)
}

// Mutation option
func Mutation(c gql.Object) schema.Option {
	return schema.Mutation(c)
}

// Subscription option
func Subscription(c gql.Object) schema.Option {
	return schema.Subscription(c)
}

// Root creates new root object
func Root() gql.Object {
	s := fmt.Sprintf("root%v", time.Now().UnixNano())
	return Object(s)
}

// Object creates new object
func Object(s string) gql.Object {
	return object.New(s)
}

// Field creates new field
func Field(s string) gql.Field {
	return field.New(s)
}

// Arg creates new argument
func Arg(s string) gql.Argument {
	return argument.New(s)
}

// List creates a output list field
func List(o interface{}) graphql.Output {
	switch v := o.(type) {
	case gql.Union:
		return graphql.NewList(v.Config().Union())
	case gql.Object:
		return graphql.NewList(v.Config().Output())
	case graphql.Type:
		return graphql.NewList(v)
	default:
		return nil
	}
}

// Multiple creates a input list field
func Multiple(o interface{}) graphql.Input {
	switch v := o.(type) {
	case gql.Object:
		return graphql.NewList(v.Config().Input())
	case graphql.Type:
		return graphql.NewList(v)
	default:
		return nil
	}
}

// Enum creates an enum
func Enum(name string) gql.Enum {
	return enum.New(name)
}

// Value creates enum value option
func Value(name string) gql.EnumValue {
	return enum.Value(name)
}

// Union creates an union
func Union(name string) gql.Union {
	return union.New(name)
}
