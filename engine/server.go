package engine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/ua-parser/uap-go/uaparser"
	"github.com/vaniila/hyper/cache"
	"github.com/vaniila/hyper/dataloader"
	"github.com/vaniila/hyper/fault"
	"github.com/vaniila/hyper/gws"
	"github.com/vaniila/hyper/logger"
	"github.com/vaniila/hyper/message"
	"github.com/vaniila/hyper/router"
	"github.com/vaniila/hyper/websocket"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"golang.org/x/net/http2"

	"net/http"
	"net/textproto"
)

type server struct {
	id         string
	addr       string
	protocol   Protocol
	cors       *cors
	cache      cache.Service
	message    message.Service
	logger     logger.Service
	gws        gws.Service
	dataloader dataloader.Service
	router     router.Service
	websocket  websocket.Service
	uaparser   *uaparser.Parser
	ln         *net.Listener
	traceid    func() string
}

func (v *server) handleParameters(c *Context, route router.RouteConfig, params []router.Param) {
	r := c.Req()
	for _, pa := range params {
		conf := pa.Config()
		data := &Value{
			typ: conf.Type(),
			key: conf.Name(),
			fmt: conf.Format(),
		}
		switch conf.Type() {
		case router.ParamBody:
			if conf.Format() == router.File {
				file, _, err := r.FormFile(conf.Name())
				if err == nil {
					b := bytes.Buffer{}
					defer file.Close()
					defer b.Reset()
					l, e := io.Copy(&b, file)
					data.val = b.Bytes()
					data.has = l > 0 && e == nil
				}
			} else if vs := r.Form[conf.Name()]; len(vs) > 0 {
				data.val = []byte(vs[0])
				data.has = true
			} else if vs := r.PostForm[conf.Name()]; len(vs) > 0 {
				data.val = []byte(vs[0])
				data.has = true
			}
		case router.ParamParam:
			data.val = []byte(chi.URLParam(r, conf.Name()))
			data.has = true
		case router.ParamQuery:
			if queries := r.URL.Query(); queries != nil {
				if vs, ok := queries[conf.Name()]; ok && len(vs) > 0 {
					data.val = []byte(vs[0])
					data.has = true
				}
			}
		case router.ParamHeader:
			if headers := textproto.MIMEHeader(r.Header); headers != nil {
				if vs, ok := headers[conf.Name()]; ok && len(vs) > 0 {
					data.val = []byte(vs[0])
					data.has = true
				}
			}
		case router.ParamCookie:
			if cookies := r.Cookies(); cookies != nil {
				for _, c := range cookies {
					if c != nil && c.Name == conf.Name() {
						data.val = []byte(c.Value)
						data.has = true
					}
				}
			}
		case router.ParamOneOf:
			if params := conf.OneOf(); len(params) > 0 {
				var fields []int
				var offset = len(c.values)
				v.handleParameters(c, route, params)
				for i := 0; i < len(params); i++ {
					value := c.values[i+offset]
					if value.Has() {
						fields = append(fields, i)
					}
				}
				if len(fields) > 1 {
					for _, field := range fields {
						param := params[field]
						conf := param.Config()
						warning := fault.
							For(fault.Conflict).
							SetResource(conf.Type().String()).
							SetField(conf.Name())
						c.warnings = append(c.warnings, warning)
					}
				}
			}
		}
		if len(data.val) == 0 || data.val == nil {
			if conf.Require() {
				warning := fault.
					For(fault.MissingField).
					SetResource(conf.Type().String()).
					SetField(conf.Name())
				c.warnings = append(c.warnings, warning)
			}
			data.val = conf.Default()
		}
		if len(data.val) != 0 && data.val != nil {
			custom := conf.Custom()
			depson := conf.DependsOn()
			switch parsed, ok := router.Val(conf.Format(), data.val); {
			case ok && (custom == nil || (custom != nil && custom(data.val))):
				data.parsed = parsed
				if len(depson) > 0 {
					for _, dep := range depson {
						idx := route.ValueIndex(dep)
						if len(c.values)-1 >= idx {
							if val := c.values[idx]; val == nil || !val.Has() {
								conf := dep.Config()
								warning := fault.
									For(fault.MissingField).
									SetResource(conf.Type().String()).
									SetField(conf.Name())
								c.warnings = append(c.warnings, warning)
							}
						}
					}
				}
			default:
				warning := fault.
					For(fault.Invalid).
					SetResource(conf.Type().String()).
					SetField(conf.Name())
				c.warnings = append(c.warnings, warning)
			}
		}
		c.values = append(c.values, data)
	}
}

func (v *server) handlerRoute(conf router.RouteConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		pid := v.traceid()

		tracer := opentracing.GlobalTracer()

		ctx, _ := tracer.Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header),
		)

		span := tracer.StartSpan(
			fmt.Sprintf("HTTP %s %s", r.Method, r.URL.String()),
			ext.RPCServerOption(ctx),
			opentracing.Tag{Key: "machine-id", Value: v.id},
			opentracing.Tag{Key: "process-id", Value: pid},
		)
		defer span.Finish()

		span.SetTag("jaeger-debug-id", pid)

		w.Header().Set("Trace-Id", pid)

		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.URL.String())
		ext.Component.Set(span, v.String())

		r = r.WithContext(
			opentracing.ContextWithSpan(r.Context(), span),
		)

		client := &Client{
			req:      r,
			uaparser: v.uaparser,
		}

		c := &Context{
			machineID:       v.id,
			processID:       pid,
			ctx:             r.Context(),
			tracer:          tracer,
			span:            span,
			identity:        new(identity),
			req:             r,
			res:             w,
			client:          client,
			values:          make([]router.Value, 0),
			params:          conf.Params(),
			warnings:        make([]fault.Cause, 0),
			cache:           v.cache,
			message:         v.message,
			logger:          v.logger,
			gqlsubscription: v.gws.Adaptor(),
			dataloader:      v.dataloader,
			dataloaders:     v.dataloader.Instance(),
			uaparser:        v.uaparser,
			statuscode:      http.StatusOK,
		}
		c.ctx = context.WithValue(c.ctx, router.RequestContext, c)

		h := &Header{
			context: c,
		}
		s := &Cookie{
			context: c,
		}
		c.header = h
		c.cookie = s

		defer func() {
			ext.HTTPStatusCode.Set(span, uint16(c.statuscode))
			if err, ok := recover().(error); ok && err != nil && !c.IsAborted() {
				ext.Error.Set(span, true)
				span.SetTag("cause", err.Error())
				c.recover = err
				if catch := conf.Catch(); catch != nil {
					catch(c)
				} else {
					c.Error(c.recover)
				}
			}
		}()

		switch r.Method {
		case "PUT", "POST", "PATCH", "CONNECT":
			span := c.StartSpan("HTTP ParseMultipartForm")
			r.Body = http.MaxBytesReader(w, r.Body, conf.MaxMemory())
			r.ParseMultipartForm(conf.MaxMemory())
			if r.MultipartForm != nil {
				defer r.MultipartForm.RemoveAll()
			}
			span.Finish()
		}
		{
			span := c.StartSpan("HTTP ParseParameters")
			v.handleParameters(c, conf, conf.Params())
			span.Finish()
		}
		if len(c.warnings) > 0 {
			err := fault.
				New("Unprocessable Entity").
				SetStatus(http.StatusUnprocessableEntity).
				AddCause(c.warnings...)
			ext.Error.Set(span, true)
			span.SetTag("cause", err.Error())
			c.Error(err)
			return
		}
		for _, md := range conf.Middlewares() {
			if !c.IsAborted() && md != nil {
				md(c)
			}
		}
		if handler := conf.Handler(); !c.IsAborted() && handler != nil {
			handler(c)
		}
	}
}

func (v *server) buildRoutes(mux *chi.Mux, routes []router.Route) {
	for _, route := range routes {
		switch conf := route.Config(); {
		case conf.Namespace():
			r := chi.NewRouter()
			v.buildRoutes(r, conf.Routes())
			mux.Mount(conf.Pattern(), r)
			for _, alias := range conf.Aliases() {
				mux.Mount(alias, r)
			}
		case conf.Method() == "GET":
			mux.Get(conf.Pattern(), v.handlerRoute(conf))
			for _, alias := range conf.Aliases() {
				mux.Get(alias, v.handlerRoute(conf))
			}
		case conf.Method() == "HEAD":
			mux.Head(conf.Pattern(), v.handlerRoute(conf))
			for _, alias := range conf.Aliases() {
				mux.Head(alias, v.handlerRoute(conf))
			}
		case conf.Method() == "OPTIONS":
			mux.Options(conf.Pattern(), v.handlerRoute(conf))
			for _, alias := range conf.Aliases() {
				mux.Options(alias, v.handlerRoute(conf))
			}
		case conf.Method() == "POST":
			mux.Post(conf.Pattern(), v.handlerRoute(conf))
			for _, alias := range conf.Aliases() {
				mux.Post(alias, v.handlerRoute(conf))
			}
		case conf.Method() == "PUT":
			mux.Put(conf.Pattern(), v.handlerRoute(conf))
			for _, alias := range conf.Aliases() {
				mux.Put(alias, v.handlerRoute(conf))
			}
		case conf.Method() == "PATCH":
			mux.Patch(conf.Pattern(), v.handlerRoute(conf))
			for _, alias := range conf.Aliases() {
				mux.Patch(alias, v.handlerRoute(conf))
			}
		case conf.Method() == "DELETE":
			mux.Delete(conf.Pattern(), v.handlerRoute(conf))
			for _, alias := range conf.Aliases() {
				mux.Delete(alias, v.handlerRoute(conf))
			}
		}
	}
	if v.websocket != nil {
		mux.Get("/_s", v.handlerRoute(
			v.router.Get("/_s").
				Name("Websocket").
				Doc(`Websocket endpoint`).
				Summary(`Websocket endpoint`).
				ClearParams().
				ClearMiddleware().
				Handle(func(c router.Context) {
					v.websocket.Handle(c)
				}).
				Config(),
		))
	}
}

func (v *server) Start() error {

	d, err := uaparser.NewFromBytes(uas)
	if err != nil {
		return err
	}
	v.uaparser = d

	// create net listener
	ln, err := net.Listen("tcp", v.addr)
	if err != nil {
		return err
	}
	v.ln = &ln

	// create router
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.DefaultCompress)
	mux.Use(middleware.Heartbeat("/healthz"))
	mux.Use(v.cors.Handler)

	v.buildRoutes(mux, v.router.Routes())

	// create http server
	srv := &http.Server{
		Addr:    v.addr,
		Handler: mux,
	}

	// enable http 2.0 if option is enabled
	if v.protocol == HTTP2 {
		http2.ConfigureServer(srv, &http2.Server{})
	}

	go srv.Serve(*v.ln)
	return nil
}

func (v *server) Stop() error {
	if v.ln != nil {
		// close net listener
		ln := *v.ln
		return ln.Close()
	}
	return nil
}

func (v *server) String() string {
	return "Hyper::Engine"
}
