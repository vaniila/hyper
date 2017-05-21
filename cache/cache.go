package cache

import "time"

// Service interface
type Service interface {
	Start() error
	Stop() error
	Set([]byte, []byte, time.Duration) error
	Get([]byte) ([]byte, error)
	String() string
}

// New creates engine server
func New(opts ...Option) Service {
	o := newOptions(opts...)
	s := &server{
		id: o.ID,
	}
	return s
}
