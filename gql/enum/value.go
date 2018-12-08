package enum

import (
	"github.com/vaniila/hyper/gql"
)

type value struct {
	name, description string
	deprecationReason string
	value             interface{}
	conf              gql.EnumValueConfig
}

func (v *value) Is(value interface{}) gql.EnumValue {
	v.value = value
	return v
}

func (v *value) Description(s string) gql.EnumValue {
	v.description = s
	return v
}

func (v *value) Deprecation(s string) gql.EnumValue {
	v.deprecationReason = s
	return v
}

func (v *value) Config() gql.EnumValueConfig {
	if v.conf == nil {
		v.conf = &valueconfig{v}
	}
	return v.conf
}

// Value creates graphql enum value
func Value(name string) gql.EnumValue {
	value := &value{
		name: name,
	}
	return value
}
