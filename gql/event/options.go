package event

// Option func
type Option func(*Options)

// Options is the engine server options
type Options struct {
	Field   string
	Payload []byte
	Filters map[string]interface{}
	EqIDs   []int64
	NeIDs   []int64
	EqKeys  []string
	NeKeys  []string
}

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Field to set field option
func Field(s string) Option {
	return func(o *Options) {
		o.Field = s
	}
}

// Payload to set payload option
func Payload(b []byte) Option {
	return func(o *Options) {
		o.Payload = b
	}
}

// Filters to set filters option
func Filters(s map[string]interface{}) Option {
	return func(o *Options) {
		o.Filters = s
	}
}

// EqIDs to set equal ids option
func EqIDs(i []int64) Option {
	return func(o *Options) {
		o.EqIDs = i
	}
}

// NeIDs to set non equal ids option
func NeIDs(i []int64) Option {
	return func(o *Options) {
		o.NeIDs = i
	}
}

// EqKeys to set equal keys option
func EqKeys(i []string) Option {
	return func(o *Options) {
		o.EqKeys = i
	}
}

// NeKeys to set non equal keys option
func NeKeys(i []string) Option {
	return func(o *Options) {
		o.NeKeys = i
	}
}
