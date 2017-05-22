package engine

// Service interface
type Service interface {
	Start() error
	Stop() error
	String() string
}

// New creates engine server
func New(opts ...Option) Service {
	o := newOptions(opts...)
	s := &server{
		id:        o.ID,
		addr:      o.Addr,
		protocol:  o.Protocol,
		cache:     o.Cache,
		message:   o.Message,
		router:    o.Router,
		websocket: o.Websocket,
	}
	return s
}
