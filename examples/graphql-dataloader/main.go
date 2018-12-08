package main

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/dataloader"
	"github.com/vaniila/hyper/gql/graphql"
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
							return r.Context().DataLoader(loader).Load(r.Context(), id)
						}),
				),
		),
		graphql.Mutation(
			graphql.
				Object("Mutation").
				Fields(
					graphql.
						Field("createNote").
						Type(object).
						Args(
							graphql.
								Arg("content").
								Type(graphql.String).
								Require(true),
						).
						Resolve(func(r graphql.Resolver) (interface{}, error) {
							id := strconv.Itoa(int(time.Now().Unix()))
							et := &note{
								id:      id,
								content: r.MustArg("content").String(),
							}
							store = append(store, et)
							return et, nil
						}),
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
			hyper.GQLBodies...,
		).
		Handle(hyper.GraphQL(schema))

	h.Run()

}
