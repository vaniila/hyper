package engine

import (
	"context"
	"net/http"

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
	cache                router.Cache
	cookie               router.Cookie
	header               router.Header
	aborted              bool
	wrote                bool
	params               []router.Param
	values               []router.Value
	uaparser             *uaparser.Parser
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

func (v *Context) Cache() router.Cache {
	return v.cache
}

func (v *Context) Cookie() router.Cookie {
	return v.cookie
}

func (v *Context) Header() router.Header {
	return v.header
}

func (v *Context) MustParam(s string) router.Value {
	return nil
}

func (v *Context) MustQuery(s string) router.Value {
	return nil
}

func (v *Context) MustBody(s string) router.Value {
	return nil
}

func (v *Context) Param(s string) (error, router.Value) {
	var exist bool
	for _, param := range v.params {
		if s == param.Config().Name() {
			exist = true
			break
		}
	}
	if !exist {
		return nil, nil
	}
	for _, value := range v.values {
		if s == value.Key() {
			return nil, value
		}
	}
	return nil, nil
}

func (v *Context) Query(s string) (error, router.Value) {
	return nil, nil
}

func (v *Context) Body(s string) (error, router.Value) {
	return nil, nil
}

func (v *Context) File(s string) []byte {
	return nil
}

func (v *Context) Tracer() tracer.Tracer {
	return nil
}

func (v *Context) Abort() {
	v.aborted = true
}

func (v *Context) IsAborted() bool {
	return v.aborted
}

func (v *Context) Write(b []byte) router.Context {
	if !v.IsAborted() {
		v.res.Write(b)
	}
	return v
}

func (v *Context) Json(o interface{}) router.Context {
	return v
}

func (v *Context) Status(code int) router.Context {
	if !v.IsAborted() && !v.wrote {
		v.res.WriteHeader(code)
		v.wrote = true
	}
	return v
}
