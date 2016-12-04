package engine

import (
	"net"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/samuelngs/hyper/router"
	"github.com/ua-parser/uap-go/uaparser"

	"golang.org/x/net/http2"

	"net/http"
)

type server struct {
	id       string
	addr     string
	protocol Protocol
	router   router.Service
	uaparser *uaparser.Parser
	ln       *net.Listener
}

func (v *server) handler(conf router.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Context()
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
			uaparser:  v.uaparser,
		}
		for _, md := range conf.Middlewares() {
			if !c.IsAborted() {
				md(c)
			}
		}
		if !c.IsAborted() {
			handler := conf.Handler()
			handler(c)
		}
	}
}

func (v *server) buildRoutes(mux *chi.Mux, routes []router.Router) {
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
