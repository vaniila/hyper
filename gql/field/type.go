package field

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql/interfaces"
)

type cast struct {
	name, description string
	fields            []interfaces.Field
}

func (v *cast) Name(s string) interfaces.Type {
	v.name = s
	return v
}

func (v *cast) Description(s string) interfaces.Type {
	v.description = s
	return v
}

func (v *cast) Fields(fs ...interfaces.Field) interfaces.Type {
	v.fields = append(v.fields, fs...)
	return v
}

func (v *cast) Input() graphql.Input {
	return nil
}

func (v *cast) Output() graphql.Output {
	return nil
}
