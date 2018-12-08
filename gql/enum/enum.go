package enum

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql"
)

var enums = make(map[string]gql.Enum)

type enum struct {
	name, description string
	values            []gql.EnumValue
	valuesMap         map[string]struct{}
	conf              graphql.EnumConfig
}

func (v *enum) Description(s string) gql.Enum {
	v.description = s
	return v
}

func (v *enum) Values(values ...gql.EnumValue) gql.Enum {
	for _, value := range values {
		if value != nil {
			name := value.Config().Name()
			if _, ok := v.valuesMap[name]; !ok {
				v.values = append(v.values, value)
				v.valuesMap[name] = struct{}{}
			}
		}
	}
	return v
}

func (v *enum) Config() gql.EnumConfig {
	return &enumconfig{
		enum: v,
	}
}

// New creates graphql enum instance
func New(name string) gql.Enum {
	if _, ok := enums[name]; !ok {
		enums[name] = &enum{
			name:      name,
			valuesMap: make(map[string]struct{}),
		}
	}
	return enums[name]
}
