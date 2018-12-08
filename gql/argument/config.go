package argument

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql"
)

type compiled struct {
	argument    *graphql.ArgumentConfig
	inputobject *graphql.InputObjectFieldConfig
}

type argumentconfig struct {
	argument *argument
	compiled *compiled
}

func (v *argumentconfig) Name() string {
	return v.argument.name
}

func (v *argumentconfig) Description() string {
	return v.argument.description
}

func (v *argumentconfig) Type() graphql.Input {
	return v.argument.typ
}

func (v *argumentconfig) Object() gql.Object {
	return v.argument.obj
}

func (v *argumentconfig) Default() interface{} {
	return v.argument.def
}

func (v *argumentconfig) Require() bool {
	return v.argument.require
}

func (v *argumentconfig) ArgumentConfig() *graphql.ArgumentConfig {
	if v.compiled == nil {
		v.compiled = new(compiled)
	}
	if v.compiled.argument == nil {
		var typ = v.argument.typ
		if v.argument.require {
			typ = graphql.NewNonNull(typ)
		}
		v.compiled.argument = &graphql.ArgumentConfig{
			Type:         typ,
			DefaultValue: v.argument.def,
			Description:  v.argument.description,
		}
	}
	return v.compiled.argument
}

func (v *argumentconfig) InputObjectFieldConfig() *graphql.InputObjectFieldConfig {
	if v.compiled == nil {
		v.compiled = new(compiled)
	}
	if v.compiled.inputobject == nil {
		var typ = v.argument.typ
		if v.argument.require {
			typ = graphql.NewNonNull(typ)
		}
		v.compiled.inputobject = &graphql.InputObjectFieldConfig{
			Type:         typ,
			DefaultValue: v.argument.def,
			Description:  v.argument.description,
		}
	}
	return v.compiled.inputobject
}
