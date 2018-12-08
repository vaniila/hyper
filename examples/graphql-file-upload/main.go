package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/gql/graphql"
)

var store []byte

// create graphql schema
var schema = graphql.
	Schema(
		graphql.Query(
			graphql.
				Object("Query").
				Fields(
					graphql.
						Field("file").
						Type(graphql.String).
						Resolve(func(r graphql.Resolver) (interface{}, error) {
							if len(store) == 0 {
								return nil, errors.New("file has not been uploaded yet")
							}
							typ := http.DetectContentType(store)
							return fmt.Sprintf("File exists: %s type, %d bytes", typ, len(store)), nil
						}),
				),
		),
		graphql.Mutation(
			graphql.
				Object("Mutation").
				Fields(
					graphql.
						Field("upload").
						Type(graphql.String).
						Resolve(func(r graphql.Resolver) (interface{}, error) {
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

	h.Run()

}
