package router

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

type server struct {
	id                   string
	notfound, notallowed HandlerFunc
	middleware           HandlerFuncs
	routes               []Route
}

func (v *server) Start() error {
	return nil
}

func (v *server) Stop() error {
	return nil
}

func (v *server) add(pat, method string) Route {
	var middleware HandlerFuncs
	if l := len(v.middleware); l > 0 {
		middleware = append(middleware, v.middleware...)
	}
	r := &router{
		pat:        pat,
		method:     method,
		ws:         true,
		http:       true,
		memory:     defaultMaxMemory,
		middleware: middleware,
	}
	v.routes = append(v.routes, r)
	return r
}

func (v *server) Get(pat string) Route {
	return v.add(pat, "GET")
}

func (v *server) Head(pat string) Route {
	return v.add(pat, "HEAD")
}

func (v *server) Options(pat string) Route {
	return v.add(pat, "OPTIONS")
}

func (v *server) Post(pat string) Route {
	return v.add(pat, "POST")
}

func (v *server) Put(pat string) Route {
	return v.add(pat, "PUT")
}

func (v *server) Patch(pat string) Route {
	return v.add(pat, "PATCH")
}

func (v *server) Delete(pat string) Route {
	return v.add(pat, "DELETE")
}

func (v *server) Namespace(pat string) Route {
	var middleware HandlerFuncs
	if l := len(v.middleware); l > 0 {
		middleware = append(middleware, v.middleware...)
	}
	r := &router{
		pat:        pat,
		namespace:  true,
		memory:     defaultMaxMemory,
		middleware: middleware,
	}
	v.routes = append(v.routes, r)
	return r
}

func (v *server) NotFound(h HandlerFunc) Service {
	v.notfound = h
	return v
}

func (v *server) MethodNotAllowed(h HandlerFunc) Service {
	v.notallowed = h
	return v
}

func (v *server) Middleware(h ...HandlerFunc) Service {
	v.middleware = append(v.middleware, h...)
	return v
}

func (v *server) RouteNotFound() HandlerFunc {
	return v.notfound
}

func (v *server) RouteNotAllowed() HandlerFunc {
	return v.notallowed
}

func (v *server) Routes() []Route {
	return v.routes
}

func (v *server) String() string {
	return "Hyper::Route"
}
