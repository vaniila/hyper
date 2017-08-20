package schema

import (
	"log"

	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql/interfaces"
)

type schema struct {
	query, mut, sub interfaces.Object
	conf            graphql.Schema
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

func (v *schema) Compile() graphql.Schema {
	c := graphql.SchemaConfig{}
	if v.query != nil {
		c.Query = v.query.ToObject()
	}
	if v.mut != nil {
		c.Mutation = v.mut.ToObject()
	}
	if v.sub != nil {
		c.Subscription = v.sub.ToObject()
	}
	i, err := graphql.NewSchema(c)
	if err != nil {
		log.Fatal(err)
	}
	v.conf = i
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
