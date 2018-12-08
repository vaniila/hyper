package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/gql/graphql"
)

type user struct {
	id, email string
}

var users = make(map[string]*user)

// email address scalar
var userEmailType = graphql.
	Scalar("EmailAddress").
	Serialize(func(value interface{}) (interface{}, error) {
		var s string
		if v, ok := value.(*string); ok {
			s = *v
		} else {
			s = fmt.Sprintf("%v", value)
		}
		if strings.Contains(s, "@") {
			return s, nil
		}
		return nil, errors.New("invalid email address")
	}).
	ParseLiteral(func(o ast.Value) (interface{}, error) {
		if val, ok := o.(*ast.StringValue); ok {
			return val.Value, nil
		}
		return nil, nil
	})

// user object
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
			Field("email").
			Type(userEmailType).
			Resolve(func(r graphql.Resolver) (interface{}, error) {
				if p, ok := r.Source().(*user); ok {
					return p.email, nil
				}
				return nil, nil
			}),
	)

// graphql schema
var schema = graphql.
	Schema(
		graphql.Query(
			graphql.
				Object("Query").
				Fields(
					graphql.
						Field("user").
						Type(userType).
						Args(
							graphql.
								Arg("id").
								Type(graphql.ID).
								Require(true),
						).
						Resolve(func(r graphql.Resolver) (interface{}, error) {
							if o, ok := users[r.MustArg("id").String()]; ok {
								return o, nil
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
						Field("createUser").
						Type(userType).
						Args(
							graphql.
								Arg("id").
								Type(graphql.ID).
								Require(true),
							graphql.
								Arg("email").
								Type(userEmailType).
								Require(true),
						).
						Resolve(func(r graphql.Resolver) (interface{}, error) {
							user := &user{
								id:    r.MustArg("id").String(),
								email: r.MustArg("email").String(),
							}
							users[user.id] = user
							return user, nil
						}),
					graphql.
						Field("updateUser").
						Type(userType).
						Args(
							graphql.
								Arg("id").
								Type(graphql.ID).
								Require(true),
							graphql.
								Arg("email").
								Type(userEmailType).
								Require(true),
						).
						Resolve(func(r graphql.Resolver) (interface{}, error) {
							if o, ok := users[r.MustArg("id").String()]; ok {
								o.email = r.MustArg("email").String()
								return o, nil
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

	h.Run()

}
