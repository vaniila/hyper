package gws

import "sync"

type subscriptions struct {
	subs map[string]Subscription
	sync.RWMutex
}

func (v *subscriptions) Has(id string) bool {
	v.RLock()
	defer v.RUnlock()
	if _, ok := v.subs[id]; ok {
		return true
	}
	return false
}

func (v *subscriptions) Add(sub Subscription, enforce ...bool) bool {
	if sub != nil && ((len(enforce) > 0 && enforce[0]) || !v.Has(sub.ID())) {
		v.Lock()
		defer v.Unlock()
		v.subs[sub.ID()] = sub
		return true
	}
	return false
}

func (v *subscriptions) Del(sub Subscription) bool {
	if sub != nil {
		v.Lock()
		defer v.Unlock()
		if _, ok := v.subs[sub.ID()]; ok {
			delete(v.subs, sub.ID())
			return true
		}
	}
	return false
}

func (v *subscriptions) Get(id string) Subscription {
	v.RLock()
	defer v.RUnlock()
	if o, ok := v.subs[id]; ok {
		return o
	}
	return nil
}

func (v *subscriptions) List() []Subscription {
	v.RLock()
	defer v.RUnlock()
	var idx int
	var arr = make([]Subscription, len(v.subs))
	for _, o := range v.subs {
		arr[idx] = o
		idx++
	}
	return arr
}

func (v *subscriptions) Len() int {
	v.RLock()
	defer v.RUnlock()
	return len(v.subs)
}
