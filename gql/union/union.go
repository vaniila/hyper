package union

import (
	"reflect"

	"github.com/vaniila/hyper/gql"
)

type union struct {
	name, description string
	resolves          map[reflect.Type]gql.Object
	conf              gql.UnionConfig
}

func (v *union) Description(s string) gql.Union {
	v.description = s
	return v
}

func (v *union) Resolve(t interface{}, o gql.Object) gql.Union {
	rv := reflect.ValueOf(t)
	if rv.IsNil() {
		return v
	}
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	ty := rv.Type()
	if _, ok := v.resolves[ty]; !ok {
		v.resolves[ty] = o
	}
	return v
}

func (v *union) Config() gql.UnionConfig {
	if v.conf == nil {
		v.conf = &unionconfig{
			union: v,
		}
	}
	return v.conf
}

// New creates new union instance
func New(s string) gql.Union {
	return &union{name: s, resolves: make(map[reflect.Type]gql.Object)}
}
