package field

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/fault"
	"github.com/vaniila/hyper/gql/interfaces"
	"github.com/vaniila/hyper/router"
)

type resolve struct {
	context interfaces.Context
	params  graphql.ResolveParams
	values  []router.Value
}

func (v *resolve) Context() interfaces.Context {
	return v.context
}

func (v *resolve) Params() graphql.ResolveParams {
	return v.params
}

func (v *resolve) Source() interface{} {
	return v.params.Source
}

func (v *resolve) Arg(s string) (router.Value, error) {
	for _, value := range v.values {
		if value.Key() == s {
			return value, nil
		}
	}
	err := fault.
		New("Illegal Field Entity").
		SetStatus(http.StatusUnprocessableEntity).
		AddCause(
			fault.
				For(fault.UnregisteredField).
				SetField(s),
		)
	return nil, err
}

func (v *resolve) MustArg(s string) router.Value {
	o, err := v.Arg(s)
	if err != nil {
		panic(err)
	}
	return o
}
