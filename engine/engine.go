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
		id:         o.ID,
		addr:       o.Addr,
		protocol:   o.Protocol,
		cache:      o.Cache,
		message:    o.Message,
		gws:        o.GQLSubscription,
		dataloader: o.DataLoader,
		router:     o.Router,
		websocket:  o.Websocket,
		logger:     o.Logger,
		traceid:    o.TraceID,
		cors:       newCors(o),
	}
	return s
}
