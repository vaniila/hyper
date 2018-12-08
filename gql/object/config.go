package object

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql"
)

var objects = map[string]*object{}

type compiled struct {
	object      *graphql.Object
	inputobject *graphql.InputObject
}

type objectconfig struct {
	object   *object
	compiled *compiled
}

func (v *objectconfig) Name() string {
	return v.object.name
}

func (v *objectconfig) Description() string {
	return v.object.description
}

func (v *objectconfig) Fields() []gql.Field {
	return v.object.fields
}

func (v *objectconfig) Args() []gql.Argument {
	return v.object.args
}

func (v *objectconfig) Output() *graphql.Object {
	if v.compiled == nil {
		v.compiled = &compiled{}
	}
	if v.compiled.object == nil {
		var fields = graphql.Fields{}
		v.compiled.object = graphql.NewObject(graphql.ObjectConfig{
			Name:        v.object.name,
			Description: v.object.description,
			Fields:      fields,
		})
		for _, field := range v.object.fields {
			c := field.Config()
			v.compiled.object.AddFieldConfig(c.Name(), c.Field())
		}
	}
	return v.compiled.object
}

func (v *objectconfig) HasOutput() bool {
	return v.compiled != nil && v.compiled.object != nil
}

func (v *objectconfig) Input() *graphql.InputObject {
	if v.compiled == nil {
		v.compiled = &compiled{}
	}
	if v.compiled.inputobject == nil {
		args := graphql.InputObjectConfigFieldMap{}
		for _, arg := range v.object.args {
			c := arg.Config()
			args[c.Name()] = c.InputObjectFieldConfig()
		}
		v.compiled.inputobject = graphql.NewInputObject(graphql.InputObjectConfig{
			Name:        v.object.name,
			Description: v.object.description,
			Fields:      args,
		})
	}
	return v.compiled.inputobject
}

func (v *objectconfig) HasInput() bool {
	return v.compiled != nil && v.compiled.inputobject != nil
}
