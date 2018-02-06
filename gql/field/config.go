package field

import (
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql/interfaces"
	"github.com/vaniila/hyper/gql/server"
	"github.com/vaniila/hyper/router"
)

type compiled struct {
	field *graphql.Field
}

type fieldconfig struct {
	field    *field
	compiled *compiled
}

func (v *fieldconfig) Name() string {
	return v.field.name
}

func (v *fieldconfig) Description() string {
	return v.field.description
}

func (v *fieldconfig) DeprecationReason() string {
	return v.field.deprecated
}

func (v *fieldconfig) Type() graphql.Output {
	return v.field.typ
}

func (v *fieldconfig) Args() []interfaces.Argument {
	return v.field.args
}

func (v *fieldconfig) Field() *graphql.Field {
	if v.compiled == nil {
		v.compiled = &compiled{}
	}
	if v.compiled.field == nil {
		var typ graphql.Output
		switch {
		case v.field.obj != nil:
			typ = v.field.obj.Config().Output()
		default:
			typ = v.field.typ
		}
		args := graphql.FieldConfigArgument{}
		for _, arg := range v.field.args {
			c := arg.Config()
			args[c.Name()] = c.ArgumentConfig()
		}
		hdfn := func(params graphql.ResolveParams) (interface{}, error) {
			c := server.FromContext(params.Context)
			r := &resolve{
				context: c,
				params:  params,
				values:  make([]interfaces.Value, len(v.field.args)),
			}
			v.Resolve(params.Args, r.values, v.field.args)
			o, err := v.field.resolve(r)
			if err != nil {
				r.Context().GraphQLError(err)
				return nil, err
			}
			return o, nil
		}
		v.compiled.field = &graphql.Field{
			Name:              v.field.name,
			Description:       v.field.description,
			DeprecationReason: v.field.deprecated,
			Type:              typ,
			Args:              args,
			Resolve:           hdfn,
		}
	}
	return v.compiled.field
}

func (v *fieldconfig) Resolve(params map[string]interface{}, values []interfaces.Value, args []interfaces.Argument) {
	for i, arg := range args {
		conf := arg.Config().ArgumentConfig()
		data := &value{
			key: arg.Config().Name(),
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
		if k, ok := params[data.key]; ok {
			switch o := k.(type) {
			case []byte:
				data.val = o
				data.has = true
				data.parsed = o
			case string:
				data.val = []byte(o)
				data.has = true
				data.parsed = o
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
				if args := arg.Config().Object().Config().Args(); len(args) > 0 {
					arr := make([]interfaces.Value, len(args))
					v.Resolve(o, arr, args)
					data.parsed = arr
				}
			case []interface{}:
				data.val = nil
				data.has = true
				data.parsed = k
			case nil:
			default:
				data.val = nil
				data.has = true
				data.parsed = k
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
