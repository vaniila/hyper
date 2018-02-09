package main

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/dataloader"
	"github.com/vaniila/hyper/gql"
	"github.com/vaniila/hyper/gql/interfaces"
)

type note struct {
	id      string
	content string
}

// initialize notes array
var store = make([]*note, 0)

// create a new note dataloader
var loader = dataloader.BatchLoader(func(ctx context.Context, keys []interface{}) []dataloader.Result {
	return dataloader.ForEach(keys, func(key interface{}) dataloader.Result {
		if s, ok := key.(string); ok {
			for _, note := range store {
				if note.id == s {
					return dataloader.Resolve(note)
				}
			}
			return dataloader.Resolve(nil)
		}
		return dataloader.Reject(errors.New("unable to recognize `key` argument"))
	})
})

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
							return r.Context().DataLoader(loader).Load(r.Context(), id)
						}),
				),
		),
		gql.Mutation(
			gql.
				Object("Mutation").
				Fields(
					gql.
						Field("createNote").
						Type(object).
						Args(
							gql.
								Arg("content").
								Type(gql.String).
								Require(true),
						).
						Resolve(func(r interfaces.Resolver) (interface{}, error) {
							id := strconv.Itoa(int(time.Now().Unix()))
							et := &note{
								id:      id,
								content: r.MustArg("content").String(),
							}
							store = append(store, et)
							return et, nil
						}),
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
							ob, err := r.Context().DataLoader(loader).Load(r.Context(), id)
							if err != nil {
								return nil, err
							}
							if nt, ok := ob.(*note); ok {
								nt.content = r.MustArg("content").String()
								return nt, nil
							}
							return nil, nil
						}),
				),
		),
	)

func main() {

	d := dataloader.New(
		dataloader.WithLoaders(loader),
	)

	h := hyper.New(
		hyper.Addr(":4000"),
		hyper.HTTP2(),
		hyper.DataLoader(d),
	)

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
