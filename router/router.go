package router

// Service interface
type Service interface {
	Start() error
	Stop() error
	Get(string) Router
	Head(string) Router
	Options(string) Router
	Post(string) Router
	Put(string) Router
	Patch(string) Router
	Delete(string) Router
	Namespace(string) Router
	Routes() []Router
	String() string
}

// Router interface
type Router interface {
	Get(string) Router
	Head(string) Router
	Options(string) Router
	Post(string) Router
	Put(string) Router
	Patch(string) Router
	Delete(string) Router
	Namespace(string) Router
	Alias(...string) Router
	Name(string) Router
	Summary(string) Router
	Doc(string) Router
	Params(...Param) Router
	Handle(HandlerFunc) Router
	Middleware(...HandlerFunc) Router
	Websocket(bool) Router
	HTTP(bool) Router
	Models(...Model) Router
	Config() Config
}

// Config interface
type Config interface {
	Pattern() string
	Aliases() []string
	Method() string
	Name() string
	Summary() string
	Doc() string
	Namespace() bool
	Websocket() bool
	HTTP() bool
	Params() []Param
	Routes() []Router
	Handler() HandlerFunc
	Middlewares() HandlerFuncs
	Model(int) interface{}
}

// Param interface
type Param interface {
	Format(DataFormat) Param
	Summary(string) Param
	Doc(string) Param
	Default([]byte) Param
	Require(bool) Param
}

// Model interface
type Model interface {
	Code() int
	Hash() string
}

// New creates engine server
func New(opts ...Option) Service {
	o := newOptions(opts...)
	s := &server{
		id:     o.ID,
		routes: make([]Router, 0),
	}
	return s
}
