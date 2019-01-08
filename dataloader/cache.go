// Package dataloader source: https://github.com/nicksrandall/dataloader
//
// dataloader is an implimentation of facebook's dataloader in go.
// See https://github.com/facebook/dataloader for more information
package dataloader

import (
	"sync"
)

// The Cache interface. If a custom cache is provided, it must implement this interface.
type Cache interface {
	Get(interface{}) (Thunk, bool)
	Set(interface{}, Thunk)
	Delete(interface{}) bool
	Clear()
}

// InMemoryCache is an in memory implementation of Cache interface.
// this simple implementation is well suited for
// a "per-request" dataloader (i.e. one that only lives
// for the life of an http request) but it not well suited
// for long lived cached items.
type InMemoryCache struct {
	items map[interface{}]Thunk
	mu    sync.RWMutex
}

// newCache constructs a new InMemoryCache
func newCache() *InMemoryCache {
	items := make(map[interface{}]Thunk)
	return &InMemoryCache{
		items: items,
	}
}

// toHashable returns the original or hash key
func toHashable(key interface{}) interface{} {
	switch key.(type) {
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, complex64, complex128, float32, float64, bool:
		return key
	default:
		i, err := toHash(key, 0)
		if err != nil {
			return key
		}
		return i
	}
}

// Set sets the `value` at `key` in the cache
func (c *InMemoryCache) Set(key interface{}, value Thunk) {
	c.mu.Lock()
	c.items[toHashable(key)] = value
	c.mu.Unlock()
}

// Get gets the value at `key` if it exsits, returns value (or nil) and bool
// indicating of value was found
func (c *InMemoryCache) Get(key interface{}) (Thunk, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[toHashable(key)]
	if !found {
		return nil, false
	}

	return item, true
}

// Delete deletes item at `key` from cache
func (c *InMemoryCache) Delete(key interface{}) bool {
	vkey := toHashable(key)
	if _, found := c.Get(vkey); found {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.items, vkey)
		return true
	}
	return false
}

// Clear clears the entire cache
func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	c.items = map[interface{}]Thunk{}
	c.mu.Unlock()
}
