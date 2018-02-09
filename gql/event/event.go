package event

import "github.com/vaniila/hyper/router"

// event struct
type event struct {
	field   string
	payload []byte
	filters map[string]interface{}
	eqIDs   []int64
	neIDs   []int64
	eqKeys  []string
	neKeys  []string
}

// Field returns field name
func (v *event) Field() string {
	return v.field
}

// Payload returns payload data
func (v *event) Payload() []byte {
	return v.payload
}

// Filters return filter columns and values
func (v *event) Filters() map[string]interface{} {
	return v.filters
}

// EqIDs return matching condition identity ids
func (v *event) EqIDs() []int64 {
	return v.eqIDs
}

// NeIDs return not matching condition identity ids condition
func (v *event) NeIDs() []int64 {
	return v.neIDs
}

// EqKeys returns matcing condition identity keys
func (v *event) EqKeys() []string {
	return v.eqKeys
}

// NeKeys returns non matching condition identity keys
func (v *event) NeKeys() []string {
	return v.neKeys
}

// New creates engine server
func New(opts ...Option) router.GQLEvent {
	o := newOptions(opts...)
	e := &event{
		field:   o.Field,
		payload: o.Payload,
		filters: o.Filters,
		eqIDs:   o.EqIDs,
		neIDs:   o.NeIDs,
		eqKeys:  o.EqKeys,
		neKeys:  o.NeKeys,
	}
	return e
}
