package engine

import (
	"net/http"

	"github.com/samuelngs/hyper/fault"
	"github.com/samuelngs/hyper/router"
	"github.com/samuelngs/hyper/router/cookie"
)

// Cookie struct
type Cookie struct {
	context *Context
}

// Set data to response cookie
func (v *Cookie) Set(key string, val string, o ...cookie.Option) {
	opts := cookie.Parse(o...)
	http.SetCookie(v.context.res, &http.Cookie{
		Name:     key,
		Value:    val,
		Path:     opts.Path,
		Domain:   opts.Domain,
		Expires:  opts.Expires,
		MaxAge:   opts.MaxAge,
		Secure:   opts.Secure,
		HttpOnly: opts.HttpOnly,
	})
}

// Get data from request cookie
func (v *Cookie) Get(key string) (router.Value, error) {
	for _, param := range v.context.params {
		if param.Config().Name() == key && param.Config().Type() == router.ParamCookie {
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
				SetResource(router.ParamCookie.String()).
				SetField(key),
		)
	return nil, err
}

// MustGet data from request cookie
func (v *Cookie) MustGet(key string) router.Value {
	val, err := v.Get(key)
	if err != nil {
		panic(err)
	}
	return val
}
