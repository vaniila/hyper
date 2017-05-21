package cookie

import "time"

// Option func
type Option func(*Options)

// Options is the cookie options
type Options struct {
	Path    string    // optional
	Domain  string    // optional
	Expires time.Time // optional

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

// Parse to parse options
func Parse(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Path to set cookie key name
func Path(s string) Option {
	return func(o *Options) {
		o.Path = s
	}
}
