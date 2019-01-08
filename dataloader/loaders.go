package dataloader

import "log"

// Loaders type
type Loaders struct {
	loaders map[Batch]DataLoader
}

// Get returns a single dataloader
func (v *Loaders) Get(h interface{}) DataLoader {
	if h == nil {
		log.Fatal("handler cannot be 'nil'")
	}
	b, ok := h.(Batch)
	if !ok {
		log.Fatal("object is not a dataloader 'Batch' type")
	}
	l, ok := v.loaders[b]
	if !ok {
		log.Fatal("dataloader could not be found")
	}
	return l
}
