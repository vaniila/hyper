package message

// Handler type
type Handler func([]byte)

// Close type
type Close func()

// Service interface
type Service interface {
	Start() error
	Stop() error
	Emit([]byte, []byte) error
	Listen([]byte, Handler) Close
	String() string
}

// New creates message server
func New(opts ...Option) Service {
	o := newOptions(opts...)
	s := &server{
		id:       o.ID,
		handlers: make([]*handler, 0),
	}
	return s
}
