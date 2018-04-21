package main

import (
	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/gql"
	"github.com/vaniila/hyper/gql/interfaces"
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
var noteAType = gql.
	Object("NoteA").
	Fields(
		gql.
			Field("id").
			Type(gql.ID).
			Resolve(func(r interfaces.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*noteA); ok {
					return p.id, nil
				}
				return nil, nil
			}),
		gql.
			Field("content").
			Type(gql.String).
			Resolve(func(r interfaces.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*noteA); ok {
					return p.content, nil
				}
				return nil, nil
			}),
	)

// create note B object
var noteBType = gql.
	Object("NoteB").
	Fields(
		gql.
			Field("id").
			Type(gql.ID).
			Resolve(func(r interfaces.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*noteA); ok {
					return p.id, nil
				}
				return nil, nil
			}),
		gql.
			Field("content").
			Type(gql.String).
			Resolve(func(r interfaces.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*noteA); ok {
					return p.content, nil
				}
				return nil, nil
			}),
	)

var object = gql.
	Union("Note").
	Description("A note item").
	Resolve(new(noteA), noteAType).
	Resolve(new(noteB), noteBType)

var items = []interface{}{
	&noteA{id: "1", content: "a"},
	&noteB{id: "2", content: "b"},
}

// create graphql schema
var schema = gql.
	Schema(
		gql.Query(
			gql.
				Object("Query").
				Fields(
					gql.
						Field("note").
						Type(object).
						Args(
							gql.
								Arg("id").
								Type(gql.ID).
								Require(true),
						).
						Resolve(func(r interfaces.Resolver) (interface{}, error) {
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

	h.
		Router().
		Get("/graphiql/*").
		Handle(hyper.GraphiQL())

	h.Run()

}
