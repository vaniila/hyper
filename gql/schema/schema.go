package schema

import "github.com/vaniila/hyper/gql"

type schema struct {
	query, mut, sub gql.Object
	conf            gql.SchemaConfig
}

func (v *schema) Query(o gql.Object) gql.Schema {
	v.query = o
	return v
}

func (v *schema) Mutation(o gql.Object) gql.Schema {
	v.mut = o
	return v
}

func (v *schema) Subscription(o gql.Object) gql.Schema {
	v.sub = o
	return v
}

func (v *schema) Config() gql.SchemaConfig {
	if v.conf == nil {
		v.conf = &schemaconfig{
			schema:   v,
			compiled: new(compiled),
		}
	}
	return v.conf
}

// New creates a new schema
func New(opt ...Option) gql.Schema {
	opts := newOptions(opt...)
	return &schema{
		query: opts.Query,
		mut:   opts.Mutation,
		sub:   opts.Subscription,
	}
}
