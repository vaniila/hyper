package scalar

import "github.com/vaniila/hyper/gql"

type scalar struct {
	name, description string
	serialize         gql.ScalarSerializeHandler
	parseValue        gql.ScalarParseValueHandler
	parseLiteral      gql.ScalarParseLiteralHandler
	config            *config
}

func (v *scalar) Description(s string) gql.Scalar {
	v.description = s
	return v
}

func (v *scalar) Serialize(hdr gql.ScalarSerializeHandler) gql.Scalar {
	v.serialize = hdr
	return v
}

func (v *scalar) ParseValue(hdr gql.ScalarParseValueHandler) gql.Scalar {
	v.parseValue = hdr
	return v
}

func (v *scalar) ParseLiteral(hdr gql.ScalarParseLiteralHandler) gql.Scalar {
	v.parseLiteral = hdr
	return v
}

func (v *scalar) Config() gql.ScalarConfig {
	if v.config == nil {
		v.config = &config{scalar: v}
	}
	return v.config
}

// New creates a scalar instance
func New(name string) gql.Scalar {
	return &scalar{name: name}
}
