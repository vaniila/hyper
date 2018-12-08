package main

import (
	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/gql/graphql"
)

type role int

const (
	admin   role = 0
	regular      = 1
)

type user struct {
	id   string
	role role
}

// create default user
var defuser = &user{
	id:   "1",
	role: admin,
}

// create user role enum
var userRoleType = graphql.
	Enum("UserRole").
	Values(
		graphql.Value("ADMIN").Is(admin),
		graphql.Value("REGULAR").Is(regular),
	)

// create user object
var userType = graphql.
	Object("User").
	Fields(
		graphql.
			Field("id").
			Type(graphql.ID).
			Resolve(func(r graphql.Resolver) (interface{}, error) {
				if u, ok := r.Source().(*user); ok {
					return u.id, nil
				}
				return nil, nil
			}),
		graphql.
			Field("role").
			Type(userRoleType).
			Resolve(func(r graphql.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*user); ok {
					return p.role, nil
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
						Field("me").
						Type(userType).
						Resolve(func(r graphql.Resolver) (interface{}, error) {
							return defuser, nil
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
