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
	Handle(HandlerFunc) Route
	Catch(HandlerFunc) Route
	Middleware(...HandlerFunc) Route
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
	Routes() []Route
	Handler() HandlerFunc
	Catch() HandlerFunc
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
	Config() ParamConfig
}

// ParamConfig interface
type ParamConfig interface {
	Name() string
	Type() ParamType
	Format() DataFormat
	Summary() string
	Doc() string
	Default() []byte
	Require() bool
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
