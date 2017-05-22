package engine

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/samuelngs/hyper/fault"
	"github.com/samuelngs/hyper/router"
	"github.com/samuelngs/hyper/tracer"
	"github.com/ua-parser/uap-go/uaparser"
)

type Context struct {
	machineID, processID string
	ctx                  context.Context
	req                  *http.Request
	res                  http.ResponseWriter
	client               router.Client
	cache                router.CacheAdaptor
	message              router.MessageAdaptor
	cookie               router.Cookie
	header               router.Header
	aborted              bool
	wrote                bool
	params               []router.Param
	values               []router.Value
	uaparser             *uaparser.Parser
	recover              error
}

func (v Context) MachineID() string {
	return v.machineID
}

func (v *Context) ProcessID() string {
	return v.processID
}

func (v *Context) Context() context.Context {
	return v.ctx
}

func (v *Context) Req() *http.Request {
	return v.req
}

func (v *Context) Res() http.ResponseWriter {
	return v.res
}

func (v *Context) Client() router.Client {
	if v.client == nil {
		v.client = &Client{
			req:      v.req,
			uaparser: v.uaparser,
		}
	}
	return v.client
}

func (v *Context) Cache() router.CacheAdaptor {
	return v.cache
}

func (v *Context) Message() router.MessageAdaptor {
	return v.message
}

func (v *Context) Cookie() router.Cookie {
	return v.cookie
}

func (v *Context) Header() router.Header {
	return v.header
}

func (v *Context) MustParam(s string) router.Value {
	val, err := v.Param(s)
	if err != nil {
		panic(err)
	}
	return val
}

func (v *Context) MustQuery(s string) router.Value {
	val, err := v.Query(s)
	if err != nil {
		panic(err)
	}
	return val
}

func (v *Context) MustBody(s string) router.Value {
	val, err := v.Body(s)
	if err != nil {
		panic(err)
	}
	return val
}

func (v *Context) Param(s string) (router.Value, error) {
	for _, param := range v.params {
		if param.Config().Name() == s && param.Config().Type() == router.ParamParam {
			for _, value := range v.values {
				if value.Key() == s {
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
				SetResource(router.ParamParam.String()).
				SetField(s),
		)
	return nil, err
}

func (v *Context) Query(s string) (router.Value, error) {
	for _, param := range v.params {
		if param.Config().Name() == s && param.Config().Type() == router.ParamQuery {
			for _, value := range v.values {
				if value.Key() == s {
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
				SetResource(router.ParamQuery.String()).
				SetField(s),
		)
	return nil, err
}

func (v *Context) Body(s string) (router.Value, error) {
	for _, param := range v.params {
		if param.Config().Name() == s && param.Config().Type() == router.ParamBody {
			for _, value := range v.values {
				if value.Key() == s {
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
				SetResource(router.ParamBody.String()).
				SetField(s),
		)
	return nil, err
}

func (v *Context) File(s string) []byte {
	return nil
}

func (v *Context) Tracer() tracer.Tracer {
	return nil
}

func (v *Context) Recover() error {
	return v.recover
}

func (v *Context) Abort() {
	v.aborted = true
}

func (v *Context) IsAborted() bool {
	return v.aborted
}

func (v *Context) Write(b []byte) router.Context {
	if !v.IsAborted() {
		if !v.wrote {
			v.wrote = true
		}
		v.res.Write(b)
	}
	return v
}

func (v *Context) Error(e error) router.Context {
	if e != nil {
		if f, ok := fault.Is(e); ok {
			v.Status(f.Status())
			v.Write(f.Json())
		} else {
			v.Write([]byte(e.Error()))
		}
	}
	return v
}

func (v *Context) Json(o interface{}) router.Context {
	switch b, e := json.Marshal(o); {
	case e != nil:
		err := fault.
			New("Problems serializing JSON").
			SetStatus(http.StatusInternalServerError).
			AddCause(
				fault.
					For(fault.Invalid).
					SetResource("response"),
			)
		panic(err)
	default:
		v.Write(b)
	}
	return v
}

func (v *Context) Status(code int) router.Context {
	if !v.IsAborted() && !v.wrote {
		v.res.WriteHeader(code)
		v.wrote = true
	}
	return v
}
