package dataloader

import (
	"context"
	"sync"
)

// Batch type
type Batch interface {
	Handle(context.Context, []interface{}) []Result
}

// Each type
type Each func(interface{}) Result

// Service interface
type Service interface {
	Start() error
	Stop() error
	Instance() DataLoaders
}

// DataLoaders interface
type DataLoaders interface {
	Get(interface{}) DataLoader
}

// DataLoader interface
type DataLoader interface {
	Load(context.Context, interface{}) (interface{}, error)
	LoadMany(context.Context, []interface{}) ([]interface{}, []error)
	Clear(interface{})
	ClearAll()
	Prime(interface{}, interface{})
}

// Result interface
type Result interface {
	Data() interface{}
	Error() error
}

// ResultMany interface
type ResultMany interface {
	Data() []interface{}
	Errors() []error
}

// New creates dataloader server
func New(opts ...Option) Service {
	o := newOptions(opts...)
	s := &server{
		id:      o.ID,
		batches: o.Batches,
		opts:    o,
	}
	return s
}

// Resolve helper
func Resolve(data interface{}, optError ...error) Result {
	res := &Return{data: data}
	if len(optError) > 0 {
		res.err = optError[0]
	}
	return res
}

// Reject helper
func Reject(err error) Result {
	return &Return{err: err}
}

// ForEach helper
func ForEach(keys []interface{}, fn Each) []Result {

	len := len(keys)
	res := make([]Result, len)

	var wg sync.WaitGroup
	wg.Add(len)

	for i, key := range keys {
		go func(idx int, key interface{}) {
			defer wg.Done()
			res[idx] = fn(key)
		}(i, key)
	}

	wg.Wait()
	return res
}
