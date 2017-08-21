package gql

import (
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql/field"
	"github.com/vaniila/hyper/gql/interfaces"
	"github.com/vaniila/hyper/gql/object"
	"github.com/vaniila/hyper/gql/schema"
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

// Schema creates new schema
func Schema(opts ...schema.Option) graphql.Schema {
	return schema.New(opts...).Compile()
}

// Query option
func Query(c interfaces.Object) schema.Option {
	return schema.Query(c)
}

// Mutation option
func Mutation(c interfaces.Object) schema.Option {
	return schema.Mutation(c)
}

// Subscription option
func Subscription(c interfaces.Object) schema.Option {
	return schema.Subscription(c)
}

// Root creates new root object
func Root() interfaces.Object {
	s := fmt.Sprintf("root%v", time.Now().UnixNano())
	return object.New(object.Name(s))
}

// Object creates new object
func Object(s string) interfaces.Object {
	return object.New(object.Name(s))
}

// Field creates new field
func Field(s string) interfaces.Field {
	return field.NewField(s)
}

// Arg creates new argument
func Arg(s string) interfaces.Argument {
	return field.NewArgument(s)
}

// List creates a list field
func List(o interface{}) graphql.Output {
	switch v := o.(type) {
	case interfaces.Object:
		return graphql.NewList(v.ToObject())
	case graphql.Type:
		return graphql.NewList(v)
	default:
		return nil
	}
}
