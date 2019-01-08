package argument

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql/interfaces"
)

type argument struct {
	name, description string
	typ               graphql.Input
	obj               interfaces.Object
	def               interface{}
	require           bool
	conf              interfaces.ArgumentConfig
}

func (v *argument) Description(s string) interfaces.Argument {
	v.description = s
	return v
}

func (v *argument) Type(o interface{}) interfaces.Argument {
	switch t := o.(type) {
	case interfaces.Object:
		v.typ = t.Config().Input()
		v.obj = t
	case graphql.Input:
		v.typ = t
	}
	return v
}

func (v *argument) Default(o interface{}) interfaces.Argument {
	v.def = o
	return v
}

func (v *argument) Require(b bool) interfaces.Argument {
	v.require = b
	return v
}

func (v *argument) Config() interfaces.ArgumentConfig {
	if v.conf == nil {
		v.conf = &argumentconfig{
			argument: v,
			compiled: new(compiled),
		}
	}
	return v.conf
}

// New creates new argument instance
func New(s string) interfaces.Argument {
	return &argument{name: s}
}
