package argument

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql"
)

type argument struct {
	name, description string
	typ               graphql.Input
	obj               gql.Object
	def               interface{}
	require           bool
	initialized       bool
	conf              gql.ArgumentConfig
}

func (v *argument) Description(s string) gql.Argument {
	v.description = s
	return v
}

func (v *argument) Type(o interface{}) gql.Argument {
	switch t := o.(type) {
	case gql.Object:
		v.typ = t.Config().Input()
		v.obj = t
	case graphql.Input:
		v.typ = t
	}
	return v
}

func (v *argument) Default(o interface{}) gql.Argument {
	v.def = o
	return v
}

func (v *argument) Require(b bool) gql.Argument {
	v.require = b
	return v
}

func (v *argument) Init(fn gql.ArgumentInitializer) gql.Argument {
	if !v.initialized {
		fn(v)
		v.initialized = true
	}
	return v
}

func (v *argument) Config() gql.ArgumentConfig {
	if v.conf == nil {
		v.conf = &argumentconfig{
			argument: v,
			compiled: new(compiled),
		}
	}
	return v.conf
}

// New creates new argument instance
func New(s string) gql.Argument {
	return &argument{name: s}
}
