package main

import (
	"context"
	"errors"

	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/dataloader"
	"github.com/vaniila/hyper/gql"
	"github.com/vaniila/hyper/gql/interfaces"
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
var object = gql.
	Object("User").
	Fields(
		gql.
			Field("id").
			Type(gql.ID).
			Resolve(func(r interfaces.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*user); ok {
					return p.id, nil
				}
				return nil, nil
			}),
	)

var _ = object.
	RecursiveFields(
		gql.
			Field("friend").
			Args(
				gql.
					Arg("id").
					Type(gql.ID).
					Require(true),
			).
			Type(object).
			Resolve(func(r interfaces.Resolver) (interface{}, error) {
				id := r.MustArg("id").String()
				req := &request{id: id}
				return r.Context().DataLoader(loader).Load(r.Context(), req)
			}),
		gql.
			Field("friends").
			Type(gql.List(object)).
			Resolve(func(r interfaces.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*user); ok {
					return p.friends, nil
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
						Field("user").
						Type(object).
						Args(
							gql.
								Arg("id").
								Type(gql.ID).
								Require(true),
						).
						Resolve(func(r interfaces.Resolver) (interface{}, error) {
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
