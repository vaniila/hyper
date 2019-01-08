package kv

import (
	"github.com/vaniila/hyper/router"
	"golang.org/x/sync/syncmap"
)

func New() router.KV {
	return &kv{new(syncmap.Map)}
}
