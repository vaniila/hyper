package sync

type channel struct {
	namespace    Namespace
	name         string
	nsubscribers []Context
	server       *server
}

func (v *channel) Namespace() Namespace {
	return v.namespace
}

func (v *channel) Name() string {
	return v.name
}

func (v *channel) NodeSubscribers() []Context {
	return v.nsubscribers
}

func (v *channel) Has(c Context) bool {
	for _, s := range v.nsubscribers {
		if s.MachineID() == c.MachineID() && s.ProcessID() == c.ProcessID() {
			return true
		}
	}
	return false
}

func (v *channel) Subscribe(c Context) Channel {
	subscribed := v.Has(c)
	if !subscribed {
		v.nsubscribers = append(v.nsubscribers, c)
		c.Subscriptions().Add(v)
	}
	return v
}

func (v *channel) Unsubscribe(c Context) Channel {
	for i, s := range v.nsubscribers {
		if s.MachineID() == c.MachineID() && s.ProcessID() == c.ProcessID() {
			v.nsubscribers = append(v.nsubscribers[:i], v.nsubscribers[i+1:]...)
			c.Subscriptions().Del(v)
		}
	}
	return v
}

func (v *channel) Write(p *Packet, s ...*Condition) error {
	var c *Condition
	for _, i := range s {
		c = i
		break
	}
	if p == nil {
		return nil
	}
	if p.Namespace == "" {
		p.Namespace = v.namespace.Config().Namespace()
	}
	if p.Channel == "" {
		p.Channel = v.name
	}
	if p.Action == ActionUnknown {
		p.Action = ActionMessage
	}
	d := &Distribution{
		Packet:    p,
		Condition: c,
	}
	return v.server.Publish(d)
}

func (v *channel) BeforeOpen() {

}

func (v *channel) AfterClose() {
}
