package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/gql"
	"github.com/vaniila/hyper/gql/interfaces"
)

var store []byte

// create graphql schema
var schema = gql.
	Schema(
		gql.Query(
			gql.
				Object("Query").
				Fields(
					gql.
						Field("file").
						Type(gql.String).
						Resolve(func(r interfaces.Resolver) (interface{}, error) {
							if len(store) == 0 {
								return nil, errors.New("file has not been uploaded yet")
							}
							typ := http.DetectContentType(store)
							return fmt.Sprintf("File exists: %s type, %d bytes", typ, len(store)), nil
						}),
				),
		),
		gql.Mutation(
			gql.
				Object("Mutation").
				Fields(
					gql.
						Field("upload").
						Type(gql.String).
						Resolve(func(r interfaces.Resolver) (interface{}, error) {
							store = r.Context().File("file")
							if len(store) == 0 {
								return nil, errors.New("file is missing")
							}
							typ := http.DetectContentType(store)
							return fmt.Sprintf("File uploaded: %s type, %d bytes", typ, len(store)), nil
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
