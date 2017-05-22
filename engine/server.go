package engine

import (
	"net"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/samuelngs/hyper/cache"
	"github.com/samuelngs/hyper/fault"
	"github.com/samuelngs/hyper/message"
	"github.com/samuelngs/hyper/router"
	"github.com/samuelngs/hyper/websocket"
	"github.com/ua-parser/uap-go/uaparser"

	"golang.org/x/net/http2"

	"net/http"
)

type server struct {
	id        string
	addr      string
	protocol  Protocol
	cache     cache.Service
	message   message.Service
	router    router.Service
	websocket websocket.Service
	uaparser  *uaparser.Parser
	ln        *net.Listener
}

func (v *server) handler(conf router.RouteConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &Client{
			req:      r,
			uaparser: v.uaparser,
		}
		c := &Context{
			machineID: v.id,
			processID: newID(),
			ctx:       r.Context(),
			req:       r,
			res:       w,
			client:    client,
			values:    make([]router.Value, 0),
			params:    conf.Params(),
			cache:     v.cache,
			message:   v.message,
			uaparser:  v.uaparser,
		}
		h := &Header{
			context: c,
		}
		s := &Cookie{
			context: c,
		}
		c.header = h
		c.cookie = s
		defer func() {
			if err, ok := recover().(error); ok && err != nil && !c.IsAborted() {
				c.recover = err
				if catch := conf.Catch(); catch != nil {
					catch(c)
				}
			}
		}()
		var warnings []fault.Cause
		for _, pa := range conf.Params() {
			conf := pa.Config()
			data := &Value{
				typ: conf.Type(),
				key: conf.Name(),
			}
			switch conf.Type() {
			case router.ParamBody:
				data.val = []byte(r.FormValue(conf.Name()))
			case router.ParamParam:
				data.val = []byte(chi.URLParam(r, conf.Name()))
			case router.ParamQuery:
				data.val = []byte(r.URL.Query().Get(conf.Name()))
			case router.ParamHeader:
				data.val = []byte(r.Header.Get(conf.Name()))
			case router.ParamCookie:
				if cookie, err := r.Cookie(conf.Name()); err != nil {
					data.val = make([]byte, 0)
				} else {
					data.val = []byte(cookie.Value)
				}
			}
			if len(data.val) == 0 || data.val == nil {
				if conf.Require() {
					warning := fault.
						For(fault.MissingField).
						SetResource(conf.Type().String()).
						SetField(conf.Name())
					warnings = append(warnings, warning)
				}
				data.val = conf.Default()
			}
			c.values = append(c.values, data)
		}
		if len(warnings) > 0 {
			err := fault.
				New("Validation Failed").
				SetStatus(http.StatusUnprocessableEntity).
				AddCause(warnings...)
			panic(err)
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
		case conf.Method() == "GET":
			mux.Get(conf.Pattern(), v.handler(conf))
			for _, alias := range conf.Aliases() {
				mux.Get(alias, v.handler(conf))
			}
		case conf.Method() == "HEAD":
			mux.Head(conf.Pattern(), v.handler(conf))
			for _, alias := range conf.Aliases() {
				mux.Head(alias, v.handler(conf))
			}
		case conf.Method() == "OPTIONS":
			mux.Options(conf.Pattern(), v.handler(conf))
			for _, alias := range conf.Aliases() {
				mux.Options(alias, v.handler(conf))
			}
		case conf.Method() == "POST":
			mux.Post(conf.Pattern(), v.handler(conf))
			for _, alias := range conf.Aliases() {
				mux.Post(alias, v.handler(conf))
			}
		case conf.Method() == "PUT":
			mux.Put(conf.Pattern(), v.handler(conf))
			for _, alias := range conf.Aliases() {
				mux.Put(alias, v.handler(conf))
			}
		case conf.Method() == "PATCH":
			mux.Patch(conf.Pattern(), v.handler(conf))
			for _, alias := range conf.Aliases() {
				mux.Patch(alias, v.handler(conf))
			}
		case conf.Method() == "DELETE":
			mux.Delete(conf.Pattern(), v.handler(conf))
			for _, alias := range conf.Aliases() {
				mux.Delete(alias, v.handler(conf))
			}
		}
	}
	if v.websocket != nil {
		mux.Get("/_s", v.handler(
			v.router.Get("/_s").
				Name("Websocket").
				Doc(`Websocket endpoint`).
				Summary(`Websocket endpoint`).
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
