package object

import (
	"github.com/vaniila/hyper/gql/interfaces"
)

// Option func
type Option func(*Options)

// Options struct
type Options struct {
	Name        string
	Description string
	Fields      []interfaces.Field
}

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Name to set object name
func Name(s string) Option {
	return func(o *Options) {
		o.Name = s
	}
}

// Description to set object description
func Description(s string) Option {
	return func(o *Options) {
		o.Description = s
	}
}

// Fields to set object fields
func Fields(f ...interfaces.Field) Option {
	return func(o *Options) {
		o.Fields = f
	}
}
