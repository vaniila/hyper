package schema

import (
	"github.com/vaniila/hyper/gql/interfaces"
)

type schema struct {
	query, mut, sub interfaces.Object
	conf            interfaces.SchemaConfig
}

func (v *schema) Query(o interfaces.Object) interfaces.Schema {
	v.query = o
	return v
}

func (v *schema) Mutation(o interfaces.Object) interfaces.Schema {
	v.mut = o
	return v
}

func (v *schema) Subscription(o interfaces.Object) interfaces.Schema {
	v.sub = o
	return v
}

func (v *schema) Config() interfaces.SchemaConfig {
	if v.conf == nil {
		v.conf = &schemaconfig{
			schema:   v,
			compiled: new(compiled),
		}
	}
	return v.conf
}

// New creates a new schema
func New(opt ...Option) interfaces.Schema {
	opts := newOptions(opt...)
	return &schema{
		query: opts.Query,
		mut:   opts.Mutation,
		sub:   opts.Subscription,
	}
}
