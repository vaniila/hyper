package union

import (
	"reflect"

	"github.com/vaniila/hyper/gql/interfaces"
)

type union struct {
	name, description string
	resolves          map[reflect.Type]interfaces.Object
	conf              interfaces.UnionConfig
}

func (v *union) Description(s string) interfaces.Union {
	v.description = s
	return v
}

func (v *union) Resolve(t interface{}, o interfaces.Object) interfaces.Union {
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

func (v *union) Config() interfaces.UnionConfig {
	if v.conf == nil {
		v.conf = &unionconfig{
			union: v,
		}
	}
	return v.conf
}

// New creates new union instance
func New(s string) interfaces.Union {
	return &union{name: s, resolves: make(map[reflect.Type]interfaces.Object)}
}
