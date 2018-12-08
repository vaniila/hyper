package enum

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql"
)

type enumconfig struct {
	enum     *enum
	compiled *graphql.Enum
}

func (v *enumconfig) Name() string {
	return v.enum.name
}

func (v *enumconfig) Description() string {
	return v.enum.description
}

func (v *enumconfig) Values() []gql.EnumValue {
	return v.enum.values
}

func (v *enumconfig) Enum() *graphql.Enum {
	if v.compiled == nil {
		config := graphql.EnumConfig{
			Name:        v.enum.name,
			Description: v.enum.description,
			Values:      make(graphql.EnumValueConfigMap),
		}
		for _, value := range v.enum.values {
			option := value.Config()
			evc := &graphql.EnumValueConfig{
				Value:             option.Is(),
				DeprecationReason: option.Deprecation(),
				Description:       option.Description(),
			}
			if evc.Value == nil {
				evc.Value = option.Name()
			}
			config.Values[option.Name()] = evc
		}
		v.compiled = graphql.NewEnum(config)
	}
	return v.compiled
}
