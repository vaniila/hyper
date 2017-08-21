package field

import (
	"strconv"
	"time"

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
		v.typ = o.ToObject()
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

func (v *field) ResolveParameters(params map[string]interface{}, values []interfaces.Value, args []interfaces.Argument) {
	for i, arg := range args {
		name, conf := arg.ToArgumentConfig()
		data := &value{
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
		if k, ok := params[name]; ok {
			switch o := k.(type) {
			case string:
				data.val = []byte(o)
				data.has = true
			case int:
				data.val = []byte(strconv.Itoa(o))
				data.has = true
				data.parsed = o
			case float64:
				data.val = []byte(strconv.FormatFloat(o, 'E', -1, 64))
				data.has = true
				data.parsed = o
			case bool:
				data.val = []byte(strconv.FormatBool(o))
				data.has = true
				data.parsed = o
			case time.Time:
				data.val = []byte(o.String())
				data.has = true
				data.parsed = o
			case map[string]interface{}:
				data.val = nil
				data.has = true
				if args := arg.InputObject().ExportArgs(); len(args) > 0 {
					arr := make([]interfaces.Value, len(args))
					v.ResolveParameters(o, arr, args)
					data.parsed = arr
				}
			}
		}
		if !data.has {
			if o, ok := conf.DefaultValue.([]byte); ok {
				data.val = o
			}
		}
		if data.val != nil && data.parsed == nil {
			if parsed, ok := router.Val(data.fmt, data.val); ok {
				data.parsed = parsed
			}
		}
		values[i] = data
	}
}

func (v *field) Compile() *graphql.Field {
	args := graphql.FieldConfigArgument{}
	for _, arg := range v.args {
		k, v := arg.ToArgumentConfig()
		args[k] = v
	}
	hdfn := func(params graphql.ResolveParams) (interface{}, error) {
		c := server.FromContext(params.Context)
		r := &resolve{
			context: c,
			params:  params,
			values:  make([]interfaces.Value, len(v.args)),
		}
		v.ResolveParameters(params.Args, r.values, v.args)
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
