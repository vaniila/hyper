package gws

import "sync"

type tree struct {
	state map[string][]Subscription
	sync.RWMutex
}

func (v *tree) Has(sub Subscription) bool {
	v.RLock()
	defer v.RUnlock()
	if sub != nil {
		for _, field := range sub.Fields() {
			if group, ok := v.state[field]; ok {
				for _, s := range group {
					if s == sub {
						return true
					}
				}
			}
		}
	}
	return false
}

func (v *tree) Add(sub Subscription, enforce ...bool) bool {
	v.Lock()
	defer v.Unlock()
	if sub != nil && ((len(enforce) > 0 && enforce[0]) || !v.Has(sub)) {
		for _, field := range sub.Fields() {
			if _, ok := v.state[field]; ok {
				v.state[field] = append(v.state[field], sub)
			} else {
				v.state[field] = []Subscription{sub}
			}
		}
		return true
	}
	return false
}

func (v *tree) Del(sub Subscription) bool {
	v.Lock()
	defer v.Unlock()
	if sub != nil {
		var count int
		for _, field := range sub.Fields() {
			if group, ok := v.state[field]; ok {
				for i, o := range group {
					if o == sub {
						v.state[field][i] = v.state[field][len(v.state[field])-1]
						v.state[field] = v.state[field][:len(v.state[field])-1]
						count++
					}
				}
			}
		}
		return count > 0
	}
	return false
}

func (v *tree) Get(fields ...string) []Subscription {
	v.RLock()
	defer v.RUnlock()
	var count, idx int
	for _, field := range fields {
		if _, ok := v.state[field]; ok {
			count += len(v.state[field])
		}
	}
	var subs = make([]Subscription, count)
	for _, field := range fields {
		if group, ok := v.state[field]; ok {
			for _, sub := range group {
				subs[idx] = sub
				idx++
			}
		}
	}
	return subs
}
