package sync

import (
	"crypto/rand"
	"fmt"

	"github.com/samuelngs/hyper/cache"
	"github.com/samuelngs/hyper/message"
)

// Option func
type Option func(*Options)

// Options is the engine server options
type Options struct {

	// engine server unique id
	ID string

	// message broker topic
	Topic []byte

	// cache engine
	Cache cache.Service

	// message broker
	Message message.Service
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
		ID: newID(),
	}
	for _, o := range opts {
		o(&opt)
	}
	if len(opt.Topic) == 0 || opt.Topic == nil {
		opt.Topic = []byte("sync")
	}
	return opt
}

// ID to change server reference id
func ID(s string) Option {
	return func(o *Options) {
		o.ID = s
	}
}

// Topic to set message broker pub-sub topic
func Topic(s string) Option {
	return func(o *Options) {
		o.Topic = []byte(s)
	}
}

// Cache to set custom cache engine
func Cache(v cache.Service) Option {
	return func(o *Options) {
		o.Cache = v
	}
}

// Message to set custom message broker
func Message(m message.Service) Option {
	return func(o *Options) {
		o.Message = m
	}
}
