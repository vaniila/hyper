package scalar

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/vaniila/hyper/gql"
)

type config struct {
	scalar   *scalar
	compiled *graphql.Scalar
}

func (v *config) Name() string {
	return v.scalar.name
}

func (v *config) Description() string {
	return v.scalar.description
}

func (v *config) Serialize() gql.ScalarSerializeHandler {
	return v.scalar.serialize
}

func (v *config) ParseValue() gql.ScalarParseValueHandler {
	return v.scalar.parseValue
}

func (v *config) ParseLiteral() gql.ScalarParseLiteralHandler {
	return v.scalar.parseLiteral
}

func (v *config) bindSerialize(o interface{}) interface{} {
	if v.scalar.serialize != nil {
		o, err := v.scalar.serialize(o)
		if err != nil {
			panic(err)
		}
		return o
	}
	if v.scalar.parseValue != nil {
		o, err := v.scalar.parseValue(o)
		if err != nil {
			panic(err)
		}
		return o
	}
	return nil
}

func (v *config) bindParseValue(o interface{}) interface{} {
	if v.scalar.parseValue != nil {
		o, err := v.scalar.parseValue(o)
		if err != nil {
			panic(err)
		}
		return o
	}
	if v.scalar.serialize != nil {
		o, err := v.scalar.serialize(o)
		if err != nil {
			panic(err)
		}
		return o
	}
	return nil
}

func (v *config) bindParseLiteral(o ast.Value) interface{} {
	if v.scalar.parseLiteral != nil {
		o, err := v.scalar.parseLiteral(o)
		if err != nil {
			return nil
		}
		return o
	}
	return nil
}

func (v *config) Scalar() *graphql.Scalar {
	if v.compiled == nil {
		options := graphql.ScalarConfig{
			Name:         v.scalar.name,
			Description:  v.scalar.description,
			Serialize:    v.bindSerialize,
			ParseValue:   v.bindParseValue,
			ParseLiteral: v.bindParseLiteral,
		}
		v.compiled = graphql.NewScalar(options)
	}
	return v.compiled
}
