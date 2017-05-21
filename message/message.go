package message

// Service interface
type Service interface {
	Start() error
	Stop() error
	Emit([]byte, []byte) error
	Listen([]byte) (<-chan []byte, chan<- struct{}, error)
	String() string
}

// New creates message server
func New(opts ...Option) Service {
	o := newOptions(opts...)
	s := &server{
		id: o.ID,
	}
	return s
}
