package field

import (
	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql/interfaces"
	"github.com/vaniila/hyper/router"
)

type argument struct {
	name, description string
	typ               graphql.Input
	format            int
	def               interface{}
	require           bool
	compiled          *graphql.ArgumentConfig
}

func (v *argument) Name(s string) interfaces.Argument {
	v.name = s
	return v
}

func (v *argument) Description(s string) interfaces.Argument {
	v.description = s
	return v
}

func (v *argument) Type(t graphql.Input) interfaces.Argument {
	switch t {
	case graphql.Int:
		v.format = router.Int
	case graphql.Float:
		v.format = router.F64
	case graphql.String:
		v.format = router.Text
	case graphql.Boolean:
		v.format = router.Bool
	case graphql.ID:
		v.format = router.Text
	case graphql.DateTime:
		v.format = router.DateTimeRFC3339
	default:
		v.format = router.Any
	}
	v.typ = t
	return v
}

func (v *argument) Default(o []byte) interfaces.Argument {
	v.def = o
	return v
}

func (v *argument) Require(b bool) interfaces.Argument {
	v.require = b
	return v
}

func (v *argument) Compile() (string, *graphql.ArgumentConfig) {
	if v.compiled == nil {
		v.compiled = &graphql.ArgumentConfig{
			Type:         v.typ,
			DefaultValue: v.def,
			Description:  v.description,
		}
	}
	return v.name, v.compiled
}

// NewArgument creates new argument instance
func NewArgument(s string) interfaces.Argument {
	return &argument{name: s}
}
