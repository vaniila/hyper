package logger

// Service interface
type Service interface {
	Start() error
	Stop() error

	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)

	String() string
}

// Field interface
type Field interface {
	Key() string
	Value() interface{}
}

// New creates logger
func New(opts ...Option) Service {
	o := newOptions(opts...)
	l := &logger{
		id: o.ID,
	}
	return l
}

// NewField creates a logger field
func NewField(key string, value interface{}) Field {
	return &field{key: key, value: value}
}
