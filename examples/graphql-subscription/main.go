package main

import (
	"bytes"
	"encoding/gob"

	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/gql"
	"github.com/vaniila/hyper/gql/event"
	"github.com/vaniila/hyper/gql/interfaces"
)

type note struct {
	id      string
	content string
}

// initialize notes array
var notes = []*note{{
	id:      "1",
	content: "#1 default content",
}, {
	id:      "2",
	content: "#2 default content",
}, {
	id:      "3",
	content: "#3 default content",
}}

// create note object
var object = gql.
	Object("Note").
	Fields(
		gql.
			Field("id").
			Type(gql.ID).
			Resolve(func(r interfaces.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*note); ok {
					return p.id, nil
				}
				return nil, nil
			}),
		gql.
			Field("content").
			Type(gql.String).
			Resolve(func(r interfaces.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*note); ok {
					return p.content, nil
				}
				return nil, nil
			}),
	)

// create graphql schema
var schema = gql.
	Schema(
		gql.Subscription(
			gql.
				Object("Subscription").
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
							if b, ok := r.Source().([]byte); ok {

								slice := make([]string, 0)

								buf := bytes.NewReader(b)
								gob.NewDecoder(buf).Decode(&slice)

								nt := &note{slice[0], slice[1]}
								return nt, nil
							}
							return nil, nil
						}),
				),
		),
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
							for _, note := range notes {
								if note.id == id {
									return note, nil
								}
							}
							return nil, nil
						}),
				),
		),
		gql.Mutation(
			gql.
				Object("Mutation").
				Fields(
					gql.
						Field("updateNote").
						Type(object).
						Args(
							gql.
								Arg("id").
								Type(gql.ID).
								Require(true),
							gql.
								Arg("content").
								Type(gql.String).
								Require(true),
						).
						Resolve(func(r interfaces.Resolver) (interface{}, error) {
							id := r.MustArg("id").String()
							for _, note := range notes {
								if note.id == id {
									note.content = r.MustArg("content").String()
									buf := &bytes.Buffer{}
									gob.NewEncoder(buf).Encode([]string{note.id, note.content})
									slice := buf.Bytes()
									r.Context().GQLSubscription().Emit(
										event.New(
											event.Field("note"),
											event.Payload(slice),
											event.Filters(map[string]interface{}{
												"id": id,
											}),
										),
									)
									return note, nil
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

	h.Gws().Schema(schema)

	h.
		Router().
		Post("/graphql").
		Params(
			append(hyper.GQLQueries, hyper.GQLBodies...)...,
		).
		Handle(hyper.GraphQL(schema))

	h.
		Router().
		Get("/graphiql/*").
		Handle(hyper.GraphiQL())

	h.Run()

}
