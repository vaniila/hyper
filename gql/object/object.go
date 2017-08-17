package object

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql/interfaces"
)

type object struct {
	name, description string
	fields            []interfaces.Field
}

func (v *object) Name(s string) interfaces.Object {
	v.name = s
	return v
}

func (v *object) Description(s string) interfaces.Object {
	v.description = s
	return v
}

func (v *object) Fields(fs ...interfaces.Field) interfaces.Object {
	v.fields = append(v.fields, fs...)
	return v
}

func (v *object) Compile() *graphql.Object {
	fields := graphql.Fields{}
	for _, f := range v.fields {
		v := f.Compile()
		fields[v.Name] = v
	}
	c := graphql.ObjectConfig{
		Name:        v.name,
		Description: v.description,
		Fields:      fields,
	}
	return graphql.NewObject(c)
}

// New creates a new object
func New(opt ...Option) interfaces.Object {
	opts := newOptions(opt...)
	return &object{
		name:        opts.Name,
		description: opts.Description,
		fields:      opts.Fields,
	}
}
