package router

// Service interface
type Service interface {
	Start() error
	Stop() error
	Get(string) Route
	Head(string) Route
	Options(string) Route
	Post(string) Route
	Put(string) Route
	Patch(string) Route
	Delete(string) Route
	Namespace(string) Route
	NotFound(HandlerFunc) Service
	MethodNotAllowed(HandlerFunc) Service
	Params(...Param) Service
	Middleware(...HandlerFunc) Service
	RouteNotFound() HandlerFunc
	RouteNotAllowed() HandlerFunc
	Routes() []Route
	String() string
}

// Route interface
type Route interface {
	Get(string) Route
	Head(string) Route
	Options(string) Route
	Post(string) Route
	Put(string) Route
	Patch(string) Route
	Delete(string) Route
	Namespace(string) Route
	Alias(...string) Route
	Name(string) Route
	Summary(string) Route
	Doc(string) Route
	Params(...Param) Route
	ClearParams() Route
	MaxMemory(int64) Route
	Handle(HandlerFunc) Route
	Catch(HandlerFunc) Route
	Middleware(...HandlerFunc) Route
	ClearMiddleware() Route
	Websocket(bool) Route
	HTTP(bool) Route
	Models(...Model) Route
	Config() RouteConfig
}

// RouteConfig interface
type RouteConfig interface {
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
	ValueIndex(Param) int
	MaxMemory() int64
	Routes() []Route
	Handler() HandlerFunc
	Catch() HandlerFunc
	Middlewares() HandlerFuncs
	Model(int) interface{}
}

// Param interface
type Param interface {
	Custom(CustomFunc) Param
	Format(int) Param
	Summary(string) Param
	Doc(string) Param
	Default([]byte) Param
	Require(bool) Param
	DependsOn(...Param) Param
	Config() ParamConfig
}

// ParamConfig interface
type ParamConfig interface {
	Name() string
	Type() ParamType
	Custom() CustomFunc
	Format() int
	Summary() string
	Doc() string
	Default() []byte
	Require() bool
	DependsOn() []Param
	OneOf() []Param
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
		routes: make([]Route, 0),
	}
	return s
}
