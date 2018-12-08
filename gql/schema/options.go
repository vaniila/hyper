package schema

import "github.com/vaniila/hyper/gql"

// Option func
type Option func(*Options)

// Options struct
type Options struct {
	Query, Mutation, Subscription gql.Object
}

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Query to set query object
func Query(c gql.Object) Option {
	return func(o *Options) {
		o.Query = c
	}
}

// Mutation to set mutation object
func Mutation(c gql.Object) Option {
	return func(o *Options) {
		o.Mutation = c
	}
}

// Subscription to set subscription object
func Subscription(c gql.Object) Option {
	return func(o *Options) {
		o.Subscription = c
	}
}
