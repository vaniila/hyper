package field

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/fault"
	"github.com/vaniila/hyper/gql"
)

type resolve struct {
	context gql.Context
	params  graphql.ResolveParams
	values  []gql.Value
}

func (v *resolve) Context() gql.Context {
	return v.context
}

func (v *resolve) Params() graphql.ResolveParams {
	return v.params
}

func (v *resolve) Source() interface{} {
	if m, ok := v.params.Source.(map[string]interface{}); ok {
		if d, ok := m["$subscription_payload$"]; ok {
			return d
		}
	}
	return v.params.Source
}

func (v *resolve) Arg(s string) (gql.Value, error) {
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

func (v *resolve) MustArg(s string) gql.Value {
	o, err := v.Arg(s)
	if err != nil {
		panic(err)
	}
	return o
}
