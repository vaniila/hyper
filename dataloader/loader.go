// Package dataloader source: https://github.com/nicksrandall/dataloader
//
// dataloader is an implimentation of facebook's dataloader in go.
// See https://github.com/facebook/dataloader for more information
package dataloader

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

// Loader implements the dataloader.Interface.
type Loader struct {
	// the batch function to be used by this loader
	batchFn Batch

	// the maximum batch size. Set to 0 if you want it to be unbounded.
	batchCap int

	// the internal cache. This packages contains a basic cache implementation but any custom cache
	// implementation could be used as long as it implements the `Cache` interface.
	cacheLock sync.Mutex
	cache     Cache
	// should we clear the cache on each batch?
	// this would allow batching but no long term caching
	clearCacheOnBatch bool

	// count of queued up items
	count int

	// the maximum input queue size. Set to 0 if you want it to be unbounded.
	inputCap int

	// the amount of time to wait before triggering a batch
	wait time.Duration

	// lock to protect the batching operations
	batchLock sync.Mutex

	// current batcher
	curBatcher *batcher

	// used to close the sleeper of the current batcher
	endSleeper chan bool

	// used by tests to prevent logs
	silent bool
}

// Thunk is a function that will block until the value (*Result) it contins is resolved.
// After the value it contains is resolved, this function will return the result.
// This function can be called many times, much like a Promise is other languages.
// The value will only need to be resolved once so subsequent calls will return immediately.
type Thunk func() (interface{}, error)

// ThunkMany is much like the Thunk func type but it contains a list of results.
type ThunkMany func() ([]interface{}, []error)

// type used to on input channel
type batchRequest struct {
	key     interface{}
	channel chan Result
}

// NewBatchedLoader constructs a new Loader with given options.
func newBatchedLoader(batchFn Batch, opts Options) DataLoader {
	return &Loader{
		batchFn:  batchFn,
		cache:    newCache(),
		batchCap: opts.BatchCapacity,
		inputCap: opts.InputCapacity,
		wait:     opts.Wait,
		silent:   opts.Silent,
	}
}

// Load load/resolves the given key, returning a channel that will contain the value and error
func (l *Loader) Load(ctx context.Context, key interface{}) (interface{}, error) {
	c := make(chan Result, 1)
	var result struct {
		mu    sync.RWMutex
		value Result
	}

	// lock to prevent duplicate keys coming in before item has been added to cache.
	l.cacheLock.Lock()
	if v, ok := l.cache.Get(key); ok {
		defer l.cacheLock.Unlock()
		return v()
	}

	thunk := func() (interface{}, error) {
		result.mu.RLock()
		resultNotSet := result.value == nil
		result.mu.RUnlock()

		if resultNotSet {
			result.mu.Lock()
			if v, ok := <-c; ok {
				result.value = v
			}
			result.mu.Unlock()
		}
		result.mu.RLock()
		defer result.mu.RUnlock()
		return result.value.Data(), result.value.Error()
	}

	l.cache.Set(key, thunk)
	l.cacheLock.Unlock()

	// this is sent to batch fn. It contains the key and the channel to return the
	// the result on
	req := &batchRequest{key, c}

	l.batchLock.Lock()
	// start the batch window if it hasn't already started.
	if l.curBatcher == nil {
		l.curBatcher = l.newBatcher(l.silent)
		// start the current batcher batch function
		go l.curBatcher.batch(ctx)
		// start a sleeper for the current batcher
		l.endSleeper = make(chan bool)
		go l.sleeper(l.curBatcher, l.endSleeper)
	}

	l.curBatcher.input <- req

	// if we need to keep track of the count (max batch), then do so.
	if l.batchCap > 0 {
		l.count++
		// if we hit our limit, force the batch to start
		if l.count == l.batchCap {
			// end the batcher synchronously here because another call to Load
			// may concurrently happen and needs to go to a new batcher.
			l.curBatcher.end()
			// end the sleeper for the current batcher.
			// this is to stop the goroutine without waiting for the
			// sleeper timeout.
			close(l.endSleeper)
			l.reset()
		}
	}
	l.batchLock.Unlock()

	return thunk()
}

// LoadOne loads the given key, doing one request instead of batch, and returns a thunk that resolves the key.
func (l *Loader) LoadOne(ctx context.Context, key interface{}) (interface{}, error) {
	// set batch capacity to 1
	l.batchCap = 1
	return l.Load(ctx, key)
}

// LoadMany loads mulitiple keys, returning a thunk (type: ThunkMany) that will resolve the keys passed in.
func (l *Loader) LoadMany(ctx context.Context, keys []interface{}) ([]interface{}, []error) {
	length := len(keys)
	data := make([]interface{}, length)
	errors := make([]error, length)
	c := make(chan ResultMany, 1)
	wg := sync.WaitGroup{}

	wg.Add(length)
	for i := range keys {
		go func(i int) {
			defer wg.Done()
			result, err := l.Load(ctx, keys[i])
			data[i] = result
			errors[i] = err
		}(i)
	}

	go func() {
		wg.Wait()
		c <- &Returns{data, errors}
		close(c)
	}()

	var result struct {
		mu    sync.RWMutex
		value ResultMany
	}

	thunkMany := func() ([]interface{}, []error) {
		result.mu.RLock()
		resultNotSet := result.value == nil
		result.mu.RUnlock()

		if resultNotSet {
			result.mu.Lock()
			if v, ok := <-c; ok {
				result.value = v
			}
			result.mu.Unlock()
		}
		result.mu.RLock()
		defer result.mu.RUnlock()
		return result.value.Data(), result.value.Errors()
	}

	return thunkMany()
}

// Clear clears the value at `key` from the cache, it it exsits. Returs self for method chaining
func (l *Loader) Clear(key interface{}) {
	l.cacheLock.Lock()
	l.cache.Delete(key)
	l.cacheLock.Unlock()
}

// ClearAll clears the entire cache. To be used when some event results in unknown invalidations.
// Returns self for method chaining.
func (l *Loader) ClearAll() {
	l.cacheLock.Lock()
	l.cache.Clear()
	l.cacheLock.Unlock()
}

// Prime adds the provided key and value to the cache. If the key already exists, no change is made.
// Returns self for method chaining
func (l *Loader) Prime(key interface{}, value interface{}) {
	if _, ok := l.cache.Get(key); !ok {
		thunk := func() (interface{}, error) {
			return value, nil
		}
		l.cache.Set(key, thunk)
	}
}

func (l *Loader) reset() {
	l.count = 0
	l.curBatcher = nil

	if l.clearCacheOnBatch {
		l.cache.Clear()
	}
}

type batcher struct {
	input    chan *batchRequest
	batchFn  Batch
	finished bool
	silent   bool
}

// newBatcher returns a batcher for the current requests
// all the batcher methods must be protected by a global batchLock
func (l *Loader) newBatcher(silent bool) *batcher {
	return &batcher{
		input:   make(chan *batchRequest, l.inputCap),
		batchFn: l.batchFn,
		silent:  silent,
	}
}

// stop receiving input and process batch function
func (b *batcher) end() {
	if !b.finished {
		close(b.input)
		b.finished = true
	}
}

// execute the batch of all items in queue
func (b *batcher) batch(ctx context.Context) {
	var keys []interface{}
	var reqs []*batchRequest

	for item := range b.input {
		keys = append(keys, item.key)
		reqs = append(reqs, item)
	}

	var items []Result
	var panicErr interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicErr = r
				if b.silent {
					return
				}
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				log.Printf("Dataloader: Panic received in batch function:: %v\n%s", panicErr, buf)
			}
		}()
		items = b.batchFn.Handle(ctx, keys)
	}()

	if panicErr != nil {
		for _, req := range reqs {
			req.channel <- &Return{err: fmt.Errorf("Panic received in batch function: %v", panicErr)}
			close(req.channel)
		}
		return
	}

	if len(items) != len(keys) {
		err := &Return{err: fmt.Errorf(`
			The batch function supplied did not return an array of responses
			the same length as the array of keys.

			Keys:
			%v

			Values:
			%v
		`, keys, items)}

		for _, req := range reqs {
			req.channel <- err
			close(req.channel)
		}

		return
	}

	for i, req := range reqs {
		req.channel <- items[i]
		close(req.channel)
	}
}

// wait the appropriate amount of time for the provided batcher
func (l *Loader) sleeper(b *batcher, close chan bool) {
	select {
	// used by batch to close early. usually triggered by max batch size
	case <-close:
		return
	// this will move this goroutine to the back of the callstack?
	case <-time.After(l.wait):
	}

	// reset
	// this is protected by the batchLock to avoid closing the batcher input
	// channel while Load is inserting a request
	l.batchLock.Lock()
	b.end()

	// We can end here also if the batcher has already been closed and a
	// new one has been created. So reset the loader state only if the batcher
	// is the current one
	if l.curBatcher == b {
		l.reset()
	}
	l.batchLock.Unlock()
}
