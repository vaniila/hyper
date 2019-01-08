package router

import "fmt"

type config struct {
	pat           string
	aliases       []string
	name          string
	method        string
	routes        []Route
	namespace     bool
	ws            bool
	http          bool
	documentation string
	summary       string
	params        []Param
	vidx          map[string]int
	memory        int64
	middleware    HandlerFuncs
	handler       HandlerFunc
	catch         HandlerFunc
	models        []Model
}

func (v *config) Pattern() string {
	return v.pat
}

func (v *config) Aliases() []string {
	return v.aliases
}

func (v *config) Method() string {
	return v.method
}

func (v *config) Name() string {
	return v.name
}

func (v *config) Summary() string {
	return v.summary
}

func (v *config) Doc() string {
	return v.documentation
}

func (v *config) Namespace() bool {
	return v.namespace
}

func (v *config) Websocket() bool {
	return v.ws
}

func (v *config) HTTP() bool {
	return v.http
}

func (v *config) Params() []Param {
	return v.params
}

func (v *config) ValueIndex(p Param) int {
	if p != nil {
		conf := p.Config()
		name := fmt.Sprintf("%v#%v", conf.Type(), conf.Name())
		if i, ok := v.vidx[name]; ok {
			return i
		}
	}
	return -1
}

func (v *config) MaxMemory() int64 {
	return v.memory
}

func (v *config) Routes() []Route {
	return v.routes
}

func (v *config) Handler() HandlerFunc {
	return v.handler
}

func (v *config) Catch() HandlerFunc {
	return v.catch
}

func (v *config) Middlewares() HandlerFuncs {
	return v.middleware
}

func (v *config) Model(code int) interface{} {
	return nil
}
