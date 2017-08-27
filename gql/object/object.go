package object

import (
	"github.com/vaniila/hyper/gql/interfaces"
)

type object struct {
	name, description string
	args              []interfaces.Argument
	fields            []interfaces.Field
	recursiveFields   []interfaces.Field
	conf              interfaces.ObjectConfig
}

func (v *object) Description(s string) interfaces.Object {
	v.description = s
	return v
}

func (v *object) Fields(fields ...interfaces.Field) interfaces.Object {
	for _, field := range fields {
		if field != nil {
			v.fields = append(v.fields, field)
		}
	}
	return v
}

func (v *object) RecursiveFields(fields ...interfaces.Field) interfaces.Object {
	for _, field := range fields {
		if field != nil {
			v.recursiveFields = append(v.recursiveFields, field)
		}
	}
	return v
}

func (v *object) Args(args ...interfaces.Argument) interfaces.Object {
	for _, arg := range args {
		if arg != nil {
			v.args = append(v.args, arg)
		}
	}
	return v
}

func (v *object) Config() interfaces.ObjectConfig {
	if v.conf == nil {
		v.conf = &objectconfig{
			object:   v,
			compiled: &compiled{},
		}
	}
	return v.conf
}

// New creates a new object
func New(name string) interfaces.Object {
	return &object{name: name}
}
