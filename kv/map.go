package kv

import (
	"github.com/vaniila/hyper/router"
	"golang.org/x/sync/syncmap"
)

type kv struct {
	m *syncmap.Map
}

func (v *kv) Set(s string, o []byte) router.KV {
	v.m.Store(s, o)
	return v
}

func (v *kv) Get(s string) []byte {
	o, ok := v.m.Load(s)
	if ok {
		return o.([]byte)
	}
	return nil
}
