package router

type server struct {
	id     string
	routes []Route
}

func (v *server) Start() error {
	return nil
}

func (v *server) Stop() error {
	return nil
}

func (v *server) add(pat, method string) Route {
	r := &router{
		pat:    pat,
		method: method,
		ws:     true,
		http:   true,
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
	r := &router{
		pat:       pat,
		namespace: true,
	}
	v.routes = append(v.routes, r)
	return r
}

func (v *server) Routes() []Route {
	return v.routes
}

func (v *server) String() string {
	return "Hyper::Route"
}
