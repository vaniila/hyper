package union

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

type unionconfig struct {
	union    *union
	compiled *graphql.Union
}

func (v *unionconfig) Name() string {
	return v.union.name
}

func (v *unionconfig) Description() string {
	return v.union.description
}

func (v *unionconfig) Union() *graphql.Union {
	if v.compiled == nil {
		var objs = make([]*graphql.Object, len(v.union.resolves))
		var indx int
		for _, o := range v.union.resolves {
			objs[indx] = o.Config().Output()
			indx++
		}
		v.compiled = graphql.NewUnion(graphql.UnionConfig{
			Name:        v.union.name,
			Description: v.union.description,
			Types:       objs,
			ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
				rv := reflect.ValueOf(p.Value)
				if rv.IsNil() {
					return nil
				}
				for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
					rv = rv.Elem()
				}
				ty := rv.Type()
				for k, v := range v.union.resolves {
					if k == ty {
						return v.Config().Output()
					}
				}
				return nil
			},
		})
	}
	return v.compiled
}
