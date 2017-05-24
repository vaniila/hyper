package sync

import "sync"

type channels struct {
	namespace Namespace
	channels  map[string]Channel
	server    *server
	sync.RWMutex
}

func (v *channels) Namespace() Namespace {
	return v.namespace
}

func (v *channels) Has(name string) bool {
	v.RLock()
	_, ok := v.channels[name]
	v.RUnlock()
	return ok
}

func (v *channels) Get(name string) Channel {
	v.RLock()
	c, ok := v.channels[name]
	v.RUnlock()
	if ok {
		return c
	}
	return nil
}

func (v *channels) Add(name string) Channels {
	v.RLock()
	_, ok := v.channels[name]
	v.RUnlock()
	if !ok {
		c := &channel{
			namespace:    v.namespace,
			name:         name,
			nsubscribers: make([]Context, 0),
			server:       v.server,
		}
		v.Lock()
		v.channels[name] = c
		v.Unlock()
		return v
	}
	return v
}

func (v *channels) Del(name string) Channels {
	v.RLock()
	_, ok := v.channels[name]
	v.RUnlock()
	if ok {
		v.Lock()
		delete(v.channels, name)
		v.Unlock()
	}
	return v
}

func (v *channels) List() map[string]Channel {
	return v.channels
}

func (v *channels) Len() int {
	return len(v.channels)
}
