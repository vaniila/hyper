package field

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql/interfaces"
	"github.com/vaniila/hyper/gql/server"
	"github.com/vaniila/hyper/router"
)

type field struct {
	name, description, deprecated string
	typ                           graphql.Output
	args                          []interfaces.Argument
	resolve                       interfaces.ResolveHandler
}

func (v *field) Name(s string) interfaces.Field {
	v.name = s
	return v
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
		v.typ = o.Compile()
	}
	return v
}

func (v *field) Args(args ...interfaces.Argument) interfaces.Field {
	v.args = append(v.args, args...)
	return v
}

func (v *field) Resolve(h interfaces.ResolveHandler) interfaces.Field {
	v.resolve = h
	return v
}

func (v *field) Compile() *graphql.Field {
	args := graphql.FieldConfigArgument{}
	for _, arg := range v.args {
		k, v := arg.Compile()
		args[k] = v
	}
	hdfn := func(params graphql.ResolveParams) (interface{}, error) {
		c := server.FromContext(params.Context)
		r := &resolve{
			context: c,
			params:  params,
			values:  make([]router.Value, len(v.args)),
		}
		for i, arg := range v.args {
			name, conf := arg.Compile()
			data := &Value{
				key: name,
			}
			switch conf.Type {
			case graphql.Int:
				data.fmt = router.Int
			case graphql.Float:
				data.fmt = router.F64
			case graphql.String:
				data.fmt = router.Text
			case graphql.Boolean:
				data.fmt = router.Bool
			case graphql.ID:
				data.fmt = router.Text
			case graphql.DateTime:
				data.fmt = router.DateTimeRFC3339
			default:
				data.fmt = router.Any
			}
			if v, ok := params.Args[name].(string); ok {
				data.val = []byte(v)
				data.has = true
			}
			switch data.has {
			case true:
				if parsed, ok := router.Val(data.fmt, data.val); ok {
					data.parsed = parsed
				}
			case false:
				if o, ok := conf.DefaultValue.([]byte); ok {
					data.val = o
				}
			}
			r.values[i] = data
		}
		o, err := v.resolve(r)
		if err != nil {
			r.Context().GraphQLError(err)
			return nil, err
		}
		return o, nil
	}
	return &graphql.Field{
		Name:              v.name,
		Description:       v.description,
		DeprecationReason: v.deprecated,
		Type:              v.typ,
		Args:              args,
		Resolve:           hdfn,
	}
}

// NewField creates new field instance
func NewField(s string) interfaces.Field {
	return &field{name: s}
}
