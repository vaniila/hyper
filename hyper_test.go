package hyper

import (
	"testing"

	"github.com/vaniila/hyper/gql"
	"github.com/vaniila/hyper/gql/interfaces"
	"github.com/vaniila/hyper/router"
	"github.com/vaniila/hyper/sync"
)

func TestNew(t *testing.T) {

	h := New(
		Addr(":4000"),
		HTTP2(),
	)

	ws := h.Sync()

	ws.BeforeOpen(func(c sync.Context) {
		c.Identity().SetID(100)
	})

	ws.AfterClose(func(c sync.Context) {
	})

	ws.Namespace("default").
		Alias("test").
		Name("DefaultNamespace").
		Doc(`Default websocket namespace`).
		Summary(`Default websocket namespace`).
		Authorize(func(n string, c sync.Context) error {
			return nil
		}).
		Middleware(func(m []byte, n sync.Channel, c sync.Context) {
		}).
		Handle("ping", func(m []byte, n sync.Channel, c sync.Context) {
			n.Write(
				&sync.Packet{Message: []byte{49, 50, 51}},
				&sync.Condition{
					NeIDs: []int64{101},
				},
			)
		}).
		Catch(func(m []byte, n sync.Channel, c sync.Context) {
		})

	ro := h.Router()

	ro.Params(
		Header("Authorization").
			Format(Text).
			Doc(`Authorization`).
			Require(false),
	)

	ro.
		Post("/graphql").
		Params(
			append(GQLQueries, GQLBodies...)...,
		).
		Handle(GraphQL(
			gql.Schema(
				gql.Query(
					gql.
						Root().
						Fields(
							gql.
								Field("hello").
								Type(gql.String).
								Args(
									gql.
										Arg("input").
										Type(gql.String).
										Require(false),
									gql.
										Arg("test").
										Require(true).
										Type(
											gql.
												Object("HelloInput").
												Args(
													gql.
														Arg("message").
														Description("some message").
														Type(gql.String),
												),
										),
								).
								Resolve(func(r interfaces.Resolver) (interface{}, error) {
									return r.MustArg("test").In("message").String(), nil
								}),
						),
				),
			),
		))

	ro.
		Get("/graphiql/*").
		Handle(GraphiQL())

	te := ro.Namespace("/test").
		Alias("/test2").
		Params(
			Query("name1").
				Format(Text).
				Doc(`namespace doc 1`).
				Summary(`namespace summary 1`).
				Default([]byte("Sam")).
				Require(false),
		).
		Middleware(func(c router.Context) {
			c.KV().Set("hello", []byte("wow"))
		})

	ha := te.Namespace("/test").
		Alias("/test2").
		Params(
			Query("name2").
				Format(Text).
				Doc(`namespace doc 2`).
				Summary(`namespace summary 2`).
				Default([]byte("Sam")).
				Require(false),
		).
		Middleware(func(c router.Context) {
		})

	ha.Get("/").
		Alias("/test").
		Name("TestIndex").
		Doc(`Test index page`).
		Summary(`Test index page`).
		Params(
			Query("greeting").
				Format(URL).
				Doc(`The greeting message`).
				Summary(`The greeting message`).
				Default([]byte("")).
				Require(false),
			OneOf(
				Query("m1").
					Format(Text).
					Default([]byte("")).
					Require(false),
				Query("m2").
					Format(Text).
					Default([]byte("")).
					DependsOn(
						Query("greeting"),
					),
			),
		).
		Models(
			Model(StatusOK, new(string)),
		).
		Middleware(func(c router.Context) {
			c.Write([]byte(c.ProcessID()))
			c.Write([]byte(" => "))
		}).
		Handle(func(c router.Context) {
			c.Write(c.MustQuery("greeting").Val())
			c.Write(c.Header().MustGet("Authorization").Val())
			c.Write([]byte("!"))
			c.Write(c.MustQuery("m1").Val())
			c.Write(c.MustQuery("m2").Val())
			c.Write(c.KV().Get("hello"))
		})

	h.Run()

}
