package field

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql"
)

type field struct {
	name, description, deprecated string
	typ                           graphql.Output
	obj                           gql.Object
	args                          []gql.Argument
	argsMap                       map[string]struct{}
	resolve                       gql.ResolveHandler
	initialized                   bool
	conf                          gql.FieldConfig
}

func (v *field) Description(s string) gql.Field {
	v.description = s
	return v
}

func (v *field) DeprecationReason(s string) gql.Field {
	v.deprecated = s
	return v
}

func (v *field) Type(t interface{}) gql.Field {
	switch o := t.(type) {
	case graphql.Output:
		v.typ = o
	case gql.Object:
		v.typ = nil
		v.obj = o
	case gql.Union:
		v.typ = o.Config().Union()
	}
	return v
}

func (v *field) Args(args ...gql.Argument) gql.Field {
	for _, arg := range args {
		if arg != nil {
			name := arg.Config().Name()
			if _, ok := v.argsMap[name]; !ok {
				v.args = append(v.args, arg)
				v.argsMap[name] = struct{}{}
			}
		}
	}
	return v
}

func (v *field) Resolve(h gql.ResolveHandler) gql.Field {
	v.resolve = h
	return v
}

func (v *field) Init(fn gql.FieldInitializer) gql.Field {
	if !v.initialized {
		v.initialized = true
		fn(v)
	}
	return v
}

func (v *field) Config() gql.FieldConfig {
	if v.conf == nil {
		v.conf = &fieldconfig{
			field:    v,
			compiled: &compiled{},
		}
	}
	return v.conf
}

// New creates new field instance
func New(s string) gql.Field {
	return &field{
		name:    s,
		argsMap: make(map[string]struct{}),
	}
}
