package engine

import (
	"net/http"

	"github.com/vaniila/hyper/fault"
	"github.com/vaniila/hyper/router"
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
		switch {
		case param.Config().Name() == key && param.Config().Type() == router.ParamHeader:
			for _, value := range v.context.values {
				if value.Key() == key {
					return value, nil
				}
			}
		case param.Config().Type() == router.ParamOneOf:
			for _, param := range param.Config().OneOf() {
				if param.Config().Name() == key && param.Config().Type() == router.ParamHeader {
					for _, value := range v.context.values {
						if value.Key() == key {
							return value, nil
						}
					}
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
