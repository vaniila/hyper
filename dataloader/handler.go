package dataloader

import "context"

type handler struct {
	h func(context.Context, []string) []Result
}

func (v *handler) Handle(c context.Context, k []string) []Result {
	return v.h(c, k)
}

// BatchLoader creates dataloader handler
func BatchLoader(h func(context.Context, []string) []Result) Batch {
	return &handler{h}
}
