package sync

type subscriptions struct {
	channels []Channel
}

func (v *subscriptions) Has(namespace, channel string) bool {
	for _, o := range v.channels {
		if o.Namespace().Config().Namespace() == namespace && o.Name() == channel {
			return true
		}
	}
	return false
}

func (v *subscriptions) Add(c Channel) Subscriptions {
	var subscribed bool
	for _, o := range v.channels {
		if o == c {
			subscribed = true
			break
		}
	}
	if !subscribed {
		v.channels = append(v.channels, c)
	}
	return v
}

func (v *subscriptions) Del(c Channel) Subscriptions {
	for i, o := range v.channels {
		if o == c {
			v.channels = append(v.channels[:i], v.channels[i+1:]...)
		}
	}
	return v
}

func (v *subscriptions) List() []Channel {
	return v.channels
}
