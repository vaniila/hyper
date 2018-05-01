package engine

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/ua-parser/uap-go/uaparser"
	"github.com/vaniila/hyper/dataloader"
	"github.com/vaniila/hyper/fault"
	"github.com/vaniila/hyper/kv"
	"github.com/vaniila/hyper/router"
)

const (
	fieldOutput = "output"
	typeJson    = "json"
	typeProto   = "proto"
)

type Context struct {
	machineID, processID string
	ctx                  context.Context
	tracer               opentracing.Tracer
	span                 opentracing.Span
	identity             *identity
	req                  *http.Request
	res                  http.ResponseWriter
	client               router.Client
	cache                router.CacheAdaptor
	message              router.MessageAdaptor
	gqlsubscription      router.GQLSubscriptionAdaptor
	dataloader           dataloader.Service
	dataloaders          dataloader.DataLoaders
	kv                   router.KV
	cookie               router.Cookie
	header               router.Header
	aborted              bool
	wrote                bool
	statuscode           int
	params               []router.Param
	values               []router.Value
	warnings             []fault.Cause
	uaparser             *uaparser.Parser
	recover              error
}

func (v *Context) Deadline() (time.Time, bool) {
	return v.ctx.Deadline()
}

func (v *Context) Done() <-chan struct{} {
	return v.ctx.Done()
}

func (v *Context) Err() error {
	return v.ctx.Err()
}

func (v *Context) Value(key interface{}) interface{} {
	return v.ctx.Value(key)
}

func (v *Context) MachineID() string {
	return v.machineID
}

func (v *Context) ProcessID() string {
	return v.processID
}

func (v *Context) Context() context.Context {
	return v.ctx
}

func (v *Context) Identity() router.Identity {
	return v.identity
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

func (v *Context) GQLSubscription() router.GQLSubscriptionAdaptor {
	return v.gqlsubscription
}

func (v *Context) DataLoader(o interface{}) router.DataLoaderAdaptor {
	return v.dataloaders.Get(o)
}

func (v *Context) KV() router.KV {
	if v.kv == nil {
		v.kv = kv.New()
	}
	return v.kv
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

func (v *Context) MatchParameter(s string, typ router.ParamType) router.Value {
	for _, param := range v.params {
		switch {
		case param.Config().Name() == s && param.Config().Type() == typ:
			for _, value := range v.values {
				if value.Key() == s {
					return value
				}
			}
		case param.Config().Type() == router.ParamOneOf:
			for _, param := range param.Config().OneOf() {
				if param.Config().Name() == s && param.Config().Type() == typ {
					for _, value := range v.values {
						if value.Key() == s {
							return value
						}
					}
				}
			}
		}
	}
	return nil
}

func (v *Context) Param(s string) (router.Value, error) {
	if v := v.MatchParameter(s, router.ParamParam); v != nil {
		return v, nil
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
	if v := v.MatchParameter(s, router.ParamQuery); v != nil {
		return v, nil
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
	if v := v.MatchParameter(s, router.ParamBody); v != nil {
		return v, nil
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
	if v := v.MatchParameter(s, router.ParamBody); v != nil {
		return v.Val()
	}
	return nil
}

func (v *Context) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	child := v.tracer.StartSpan(
		operationName,
		append(opts, opentracing.ChildOf(v.span.Context()))...,
	)
	return child
}

func (v *Context) Tracer() opentracing.Tracer {
	return v.tracer
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

func (v *Context) Proto(i router.ProtoMessage) router.Context {
	var (
		typ = typeJson
		out []byte
	)
	for _, value := range v.values {
		if value.Key() == fieldOutput {
			switch value.String() {
			case typeJson:
				typ = typeJson
			case typeProto:
				typ = typeProto
			}
		}
	}
	switch typ {
	case typeJson:
		out, _ = json.Marshal(i)
	case typeProto:
		out, _ = proto.Marshal(i)
	}
	return v.Write(out)
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
		ext.Error.Set(v.span, true)
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
		v.statuscode = code
		v.wrote = true
	}
	return v
}

func (v *Context) Child() router.Context {
	child := new(Context)
	*child = *v

	child.dataloaders = v.dataloader.Instance()
	child.ctx = context.WithValue(child.ctx, router.RequestContext, child)

	return child
}
