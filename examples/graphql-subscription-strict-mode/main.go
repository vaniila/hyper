package main

import (
	"bytes"
	"encoding/gob"

	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/gql/event"
	"github.com/vaniila/hyper/gql/graphql"
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
var object = graphql.
	Object("Note").
	Fields(
		graphql.
			Field("id").
			Type(graphql.ID).
			Resolve(func(r graphql.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*note); ok {
					return p.id, nil
				}
				return nil, nil
			}),
		graphql.
			Field("content").
			Type(graphql.String).
			Resolve(func(r graphql.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*note); ok {
					return p.content, nil
				}
				return nil, nil
			}),
	)

// create graphql schema
var schema = graphql.
	Schema(
		graphql.Subscription(
			graphql.
				Object("Subscription").
				Fields(
					graphql.
						Field("noteUpdated").
						Type(object).
						Args(
							graphql.
								Arg("id").
								Type(graphql.ID).
								Require(false),
						).
						Resolve(func(r graphql.Resolver) (interface{}, error) {
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
							for _, note := range notes {
								if note.id == id {
									return note, nil
								}
							}
							return nil, nil
						}),
				),
		),
		graphql.Mutation(
			graphql.
				Object("Mutation").
				Fields(
					graphql.
						Field("updateNote").
						Type(object).
						Args(
							graphql.
								Arg("id").
								Type(graphql.ID).
								Require(true),
							graphql.
								Arg("content").
								Type(graphql.String).
								Require(true),
						).
						Resolve(func(r graphql.Resolver) (interface{}, error) {
							id := r.MustArg("id").String()
							for _, note := range notes {
								if note.id == id {
									note.content = r.MustArg("content").String()
									buf := &bytes.Buffer{}
									gob.NewEncoder(buf).Encode([]string{note.id, note.content})
									slice := buf.Bytes()
									r.Context().GQLSubscription().Emit(
										event.New(
											event.Field("noteUpdated"),
											event.Payload(slice),
											event.Filters(map[string]interface{}{
												"id": id,
											}),
											event.Strict(true),
										),
									)
									r.Context().GQLSubscription().Emit(
										event.New(
											event.Field("noteUpdated"),
											event.Payload(slice),
											event.Strict(true),
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
			hyper.GQLBodies...,
		).
		Handle(hyper.GraphQL(schema))

	h.Run()

}
