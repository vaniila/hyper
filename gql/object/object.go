package object

import "github.com/vaniila/hyper/gql"

type object struct {
	name, description string
	args              []gql.Argument
	argsMap           map[string]struct{}
	fields            []gql.Field
	fieldsMap         map[string]struct{}
	conf              gql.ObjectConfig
}

func (v *object) Description(s string) gql.Object {
	v.description = s
	return v
}

func (v *object) Fields(fields ...gql.Field) gql.Object {
	for _, field := range fields {
		if field != nil {
			name := field.Config().Name()
			if _, ok := v.fieldsMap[name]; !ok {
				v.fields = append(v.fields, field)
				v.fieldsMap[name] = struct{}{}
				if v.conf != nil && v.conf.HasOutput() {
					v.conf.Output().AddFieldConfig(field.Config().Name(), field.Config().Field())
				}
			}
		}
	}
	return v
}

func (v *object) Args(args ...gql.Argument) gql.Object {
	for _, arg := range args {
		if arg != nil {
			name := arg.Config().Name()
			if _, ok := v.argsMap[name]; !ok {
				v.args = append(v.args, arg)
				v.argsMap[name] = struct{}{}
			}
		}
	}
	return v
}

func (v *object) Config() gql.ObjectConfig {
	if v.conf == nil {
		v.conf = &objectconfig{
			object:   v,
			compiled: &compiled{},
		}
	}
	return v.conf
}

// New creates a new object
func New(name string) gql.Object {
	if _, ok := objects[name]; !ok {
		objects[name] = &object{
			name:      name,
			argsMap:   make(map[string]struct{}),
			fieldsMap: make(map[string]struct{}),
		}
	}
	return objects[name]
}
