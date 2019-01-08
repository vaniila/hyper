package dataloader

import (
	"crypto/rand"
	"fmt"
	"time"
)

// Option func
type Option func(*Options)

// Options is the engine server options
type Options struct {
	// server unique id
	ID string
	// dataloader batch handlers
	Batches []Batch
	// the maximum batch size. Set to 0 if you want it to be unbounded.
	BatchCapacity int
	// the maximum input queue size. Set to 0 if you want it to be unbounded.
	InputCapacity int
	// the amount of time to wait before triggering a batch
	Wait time.Duration
	// used by tests to prevent logs
	Silent bool
}

func newID() string {
	b := new([16]byte)
	rand.Read(b[:])
	b[8] = (b[8] | 0x40) & 0x7F
	b[6] = (b[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func newOptions(opts ...Option) Options {
	opt := Options{
		ID:            newID(),
		InputCapacity: 1000,
		Wait:          16 * time.Millisecond,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// ID to change server reference id
func ID(s string) Option {
	return func(o *Options) {
		o.ID = s
	}
}

// WithLoaders to add dataloaders
func WithLoaders(batches ...Batch) Option {
	return func(o *Options) {
		for _, batch := range batches {
			if batch != nil {
				o.Batches = append(o.Batches, batch)
			}
		}
	}
}

// WithBatchCapacity sets the batch capacity. Default is 0 (unbounded).
func WithBatchCapacity(c int) Option {
	return func(o *Options) {
		o.BatchCapacity = c
	}
}

// WithInputCapacity sets the input capacity. Default is 1000.
func WithInputCapacity(c int) Option {
	return func(o *Options) {
		o.InputCapacity = c
	}
}

// WithWait sets the amount of time to wait before triggering a batch.
// Default duration is 16 milliseconds.
func WithWait(d time.Duration) Option {
	return func(o *Options) {
		o.Wait = d
	}
}

// withSilentLogger turns of log messages. It's used by the tests
func withSilentLogger() Option {
	return func(o *Options) {
		o.Silent = true
	}
}
