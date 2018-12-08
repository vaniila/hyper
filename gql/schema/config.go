package schema

import (
	"log"

	"github.com/graphql-go/graphql"
	"github.com/vaniila/hyper/gql"
)

type compiled struct {
	schema *graphql.Schema
}

type schemaconfig struct {
	schema   *schema
	compiled *compiled
}

func (v *schemaconfig) Query() gql.Object {
	return v.schema.query
}

func (v *schemaconfig) Mutation() gql.Object {
	return v.schema.mut
}

func (v *schemaconfig) Subscription() gql.Object {
	return v.schema.sub
}

func (v *schemaconfig) Schema() graphql.Schema {
	if v.compiled == nil {
		v.compiled = new(compiled)
	}
	if v.compiled.schema == nil {
		c := graphql.SchemaConfig{}
		if v.schema.query != nil {
			c.Query = v.schema.query.Config().Output()
		}
		if v.schema.mut != nil {
			c.Mutation = v.schema.mut.Config().Output()
		}
		if v.schema.sub != nil {
			c.Subscription = v.schema.sub.Config().Output()
		}
		i, err := graphql.NewSchema(c)
		if err != nil {
			log.Fatal(err)
		}
		v.compiled.schema = &i
	}
	return *v.compiled.schema
}
