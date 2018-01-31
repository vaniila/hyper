package dataloader

import "context"

type handler struct {
	h func(context.Context, []interface{}) []Result
}

func (v *handler) Handle(c context.Context, k []interface{}) []Result {
	return v.h(c, k)
}

// BatchLoader creates dataloader handler
func BatchLoader(h func(context.Context, []interface{}) []Result) Batch {
	return &handler{h}
}
