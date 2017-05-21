package cache

import (
	"time"

	"github.com/samuelngs/hyper/cache/default"
)

type server struct {
	id    string
	cache *builtin.Cache
}

func (v *server) Start() error {
	v.cache = builtin.New(builtin.NoExpiration, 10*time.Minute)
	return nil
}

func (v *server) Stop() error {
	v.cache = nil
	return nil
}

func (v *server) Set(key, val []byte, ttl time.Duration) error {
	v.cache.Set(string(key[:]), val, ttl)
	return nil
}

func (v *server) Get(key []byte) ([]byte, error) {
	blob, found := v.cache.Get(string(key[:]))
	if !found {
		return nil, nil
	}
	bytes, ok := blob.([]byte)
	if ok {
		return bytes, nil
	}
	return nil, nil
}

func (v *server) String() string {
	return "Hyper::Cache"
}
