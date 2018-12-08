package main

import (
	"context"
	"errors"

	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/dataloader"
	"github.com/vaniila/hyper/gql"
	"github.com/vaniila/hyper/gql/graphql"
)

type user struct {
	id      string
	friends []*user
}

type request struct {
	id string
}

var userA = &user{id: "1"}
var userB = &user{id: "2"}
var users = func() []*user {
	var l = []*user{userA, userB}
	l[0].friends = l
	l[1].friends = l
	return l
}()

// create a new user dataloader
var loader = dataloader.BatchLoader(func(ctx context.Context, keys []interface{}) []dataloader.Result {
	return dataloader.ForEach(keys, func(key interface{}) dataloader.Result {
		if s, ok := key.(*request); ok {
			for _, user := range users {
				if user.id == s.id {
					return dataloader.Resolve(user)
				}
			}
			return dataloader.Resolve(nil)
		}
		return dataloader.Reject(errors.New("unable to recognize `key` argument"))
	})
})

// create user object
var object = graphql.
	Object("User").
	Init(func(object gql.Object) {
		object.Fields(
			graphql.
				Field("id").
				Type(graphql.ID).
				Resolve(func(r graphql.Resolver) (interface{}, error) {
					if p, ok := r.Source().(*user); ok {
						return p.id, nil
					}
					return nil, nil
				}),
			graphql.
				Field("friend").
				Args(
					graphql.
						Arg("id").
						Type(graphql.ID).
						Require(true),
				).
				Type(object).
				Resolve(func(r graphql.Resolver) (interface{}, error) {
					id := r.MustArg("id").String()
					req := &request{id: id}
					return r.Context().DataLoader(loader).Load(r.Context(), req)
				}),
			graphql.
				Field("friends").
				Type(graphql.List(object)).
				Resolve(func(r graphql.Resolver) (interface{}, error) {
					if p, ok := r.Source().(*user); ok {
						return p.friends, nil
					}
					return nil, nil
				}),
		)
	})

// create graphql schema
var schema = graphql.
	Schema(
		graphql.Query(
			graphql.
				Object("Query").
				Fields(
					graphql.
						Field("user").
						Type(object).
						Args(
							graphql.
								Arg("id").
								Type(graphql.ID).
								Require(true),
						).
						Resolve(func(r graphql.Resolver) (interface{}, error) {
							id := r.MustArg("id").String()
							req := &request{id: id}
							return r.Context().DataLoader(loader).Load(r.Context(), req)
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
