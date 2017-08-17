package schema

import "github.com/vaniila/hyper/gql/interfaces"

// Option func
type Option func(*Options)

// Options struct
type Options struct {
	Query, Mutation, Subscription interfaces.Object
}

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Query to set query object
func Query(c interfaces.Object) Option {
	return func(o *Options) {
		o.Query = c
	}
}

// Mutation to set mutation object
func Mutation(c interfaces.Object) Option {
	return func(o *Options) {
		o.Mutation = c
	}
}

// Subscription to set subscription object
func Subscription(c interfaces.Object) Option {
	return func(o *Options) {
		o.Subscription = c
	}
}
