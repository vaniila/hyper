package main

import (
	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/gql/graphql"
)

type noteA struct {
	id      string
	content string
}

type noteB struct {
	id      string
	content string
}

// create note A object
var noteAType = graphql.
	Object("NoteA").
	Fields(
		graphql.
			Field("id").
			Type(graphql.ID).
			Resolve(func(r graphql.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*noteA); ok {
					return p.id, nil
				}
				return nil, nil
			}),
		graphql.
			Field("content").
			Type(graphql.String).
			Resolve(func(r graphql.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*noteA); ok {
					return p.content, nil
				}
				return nil, nil
			}),
	)

// create note B object
var noteBType = graphql.
	Object("NoteB").
	Fields(
		graphql.
			Field("id").
			Type(graphql.ID).
			Resolve(func(r graphql.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*noteA); ok {
					return p.id, nil
				}
				return nil, nil
			}),
		graphql.
			Field("content").
			Type(graphql.String).
			Resolve(func(r graphql.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*noteA); ok {
					return p.content, nil
				}
				return nil, nil
			}),
	)

var object = graphql.
	Union("Note").
	Description("A note item").
	Resolve(new(noteA), noteAType).
	Resolve(new(noteB), noteBType)

var items = []interface{}{
	&noteA{id: "1", content: "a"},
	&noteB{id: "2", content: "b"},
}

// create graphql schema
var schema = graphql.
	Schema(
		graphql.Query(
			graphql.
				Object("Query").
				Fields(
					graphql.
						Field("note").
						Type(object).
						Args(
							graphql.
								Arg("id").
								Type(graphql.ID).
								Require(true),
						).
						Resolve(func(r graphql.Resolver) (interface{}, error) {
							id := r.MustArg("id").String()
							for _, note := range items {
								if o, ok := note.(*noteA); ok {
									if o.id == id {
										return o, nil
									}
								}
								if o, ok := note.(*noteB); ok {
									if o.id == id {
										return o, nil
									}
								}
							}
							return nil, nil
						}),
				),
		),
	)

func main() {

	h := hyper.New(
		hyper.Addr(":4000"),
		hyper.HTTP2(),
	)

	h.
		Router().
		Post("/graphql").
		Params(
			hyper.GQLBodies...,
		).
		Handle(hyper.GraphQL(schema))

	h.Run()

}
