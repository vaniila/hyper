package engine

import (
	"net/http"

	"github.com/samuelngs/hyper/fault"
	"github.com/samuelngs/hyper/router"
)

// Header struct
type Header struct {
	context *Context
}

// Set data to response header
func (v *Header) Set(key string, val string) {
	v.context.res.Header().Set(key, val)
}

// Get data from request header
func (v *Header) Get(key string) (router.Value, error) {
	for _, param := range v.context.params {
		if param.Config().Name() == key && param.Config().Type() == router.ParamParam {
			for _, value := range v.context.values {
				if value.Key() == key {
					return value, nil
				}
			}
		}
	}
	err := fault.
		New("Illegal Field Entity").
		SetStatus(http.StatusUnprocessableEntity).
		AddCause(
			fault.
				For(fault.UnregisteredField).
				SetResource(router.ParamHeader.String()).
				SetField(key),
		)
	return nil, err
}

// MustGet data from request header
func (v *Header) MustGet(key string) router.Value {
	val, err := v.Get(key)
	if err != nil {
		panic(err)
	}
	return val
}
