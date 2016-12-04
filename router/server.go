package router

type server struct {
	id     string
	routes []Router
}

func (v *server) Start() error {
	return nil
}

func (v *server) Stop() error {
	return nil
}

func (v *server) add(pat, method string) Router {
	r := &router{
		pat:    pat,
		method: method,
		ws:     true,
		http:   true,
	}
	v.routes = append(v.routes, r)
	return r
}

func (v *server) Get(pat string) Router {
	return v.add(pat, "GET")
}

func (v *server) Head(pat string) Router {
	return v.add(pat, "HEAD")
}

func (v *server) Options(pat string) Router {
	return v.add(pat, "OPTIONS")
}

func (v *server) Post(pat string) Router {
	return v.add(pat, "POST")
}

func (v *server) Put(pat string) Router {
	return v.add(pat, "PUT")
}

func (v *server) Patch(pat string) Router {
	return v.add(pat, "PATCH")
}

func (v *server) Delete(pat string) Router {
	return v.add(pat, "DELETE")
}

func (v *server) Namespace(pat string) Router {
	r := &router{
		pat:       pat,
		namespace: true,
	}
	v.routes = append(v.routes, r)
	return r
}

func (v *server) Routes() []Router {
	return v.routes
}

func (v *server) String() string {
	return "Hyper::Router"
}
