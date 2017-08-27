package field

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql/interfaces"
)

type field struct {
	name, description, deprecated string
	typ                           graphql.Output
	obj                           interfaces.Object
	args                          []interfaces.Argument
	resolve                       interfaces.ResolveHandler
	conf                          interfaces.FieldConfig
}

func (v *field) Description(s string) interfaces.Field {
	v.description = s
	return v
}

func (v *field) DeprecationReason(s string) interfaces.Field {
	v.deprecated = s
	return v
}

func (v *field) Type(t interface{}) interfaces.Field {
	switch o := t.(type) {
	case graphql.Output:
		v.typ = o
	case interfaces.Object:
		v.typ = nil
		v.obj = o
	}
	return v
}

func (v *field) Args(args ...interfaces.Argument) interfaces.Field {
	for _, arg := range args {
		if arg != nil {
			v.args = append(v.args, arg)
		}
	}
	return v
}

func (v *field) Resolve(h interfaces.ResolveHandler) interfaces.Field {
	v.resolve = h
	return v
}

func (v *field) Config() interfaces.FieldConfig {
	if v.conf == nil {
		v.conf = &fieldconfig{
			field:    v,
			compiled: &compiled{},
		}
	}
	return v.conf
}

// New creates new field instance
func New(s string) interfaces.Field {
	return &field{name: s}
}
