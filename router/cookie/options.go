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

// Domain to set cookie domain
func Domain(s string) Option {
	return func(o *Options) {
		o.Domain = s
	}
}

// Expires to set cookie expire date
func Expires(t time.Time) Option {
	return func(o *Options) {
		o.Expires = t
	}
}

// MaxAge to set cookie maxage
func MaxAge(i int) Option {
	return func(o *Options) {
		o.MaxAge = i
	}
}

// Secure to set cookie secure
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// HttpOnly to set cookie http only
func HttpOnly(b bool) Option {
	return func(o *Options) {
		o.HttpOnly = b
	}
}
