package hyper

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/vaniila/hyper/dataloader"
	"github.com/vaniila/hyper/gql"
	"github.com/vaniila/hyper/gql/event"
	"github.com/vaniila/hyper/gql/interfaces"
	"github.com/vaniila/hyper/router"
	"github.com/vaniila/hyper/sync"
)

type note struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

var notes = []*note{
	{
		ID:      "1",
		Content: "default",
	},
}

type person struct {
	id int
}

type dataA struct {
	a string
}

type dataB struct {
	b string
}

var dataArray = []interface{}{
	&dataA{a: "data-1-a"},
	&dataA{a: "data-2-a"},
	&dataB{b: "data-3-b"},
	&dataA{a: "data-4-a"},
	&dataB{b: "data-5-b"},
}

func TestNew(t *testing.T) {

	var a = dataloader.BatchLoader(func(ctx context.Context, keys []interface{}) []dataloader.Result {
		return dataloader.ForEach(keys, func(key interface{}) dataloader.Result {
			return dataloader.Resolve(2)
		})
	})

	d := dataloader.New(
		dataloader.WithLoaders(a),
	)

	h := New(
		Addr(":4000"),
		HTTP2(),
		DataLoader(d),
	)

	ws := h.Sync()
	gw := h.Gws()
	ro := h.Router()

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

	ro.Params(
		Header("Authorization").
			Format(Text).
			Doc(`Authorization`).
			Require(false),
	)

	nt := gql.
		Object("Note").
		Fields(
			gql.
				Field("id").
				Type(gql.ID).
				Resolve(func(r interfaces.Resolver) (interface{}, error) {
					if p, ok := r.Source().(*note); ok {
						return p.ID, nil
					}
					return nil, nil
				}),
			gql.
				Field("content").
				Type(gql.String).
				Resolve(func(r interfaces.Resolver) (interface{}, error) {
					if p, ok := r.Source().(*note); ok {
						return p.Content, nil
					}
					return nil, nil
				}),
		)

	da := gql.
		Object("DataA").
		Fields(
			gql.
				Field("a").
				Type(gql.String).
				Resolve(func(r interfaces.Resolver) (interface{}, error) {
					if p, ok := r.Source().(*dataA); ok {
						return p.a, nil
					}
					return nil, nil
				}),
		)

	db := gql.
		Object("DataB").
		Fields(
			gql.
				Field("b").
				Type(gql.String).
				Resolve(func(r interfaces.Resolver) (interface{}, error) {
					if p, ok := r.Source().(*dataB); ok {
						return p.b, nil
					}
					return nil, nil
				}),
		)

	dd := gql.
		Union("Data").
		Resolve(&dataA{}, da).
		Resolve(new(dataB), db)

	ur := gql.
		Object("User").
		Fields(
			gql.
				Field("id").
				Type(gql.Int).
				Resolve(func(r interfaces.Resolver) (interface{}, error) {
					if p, ok := r.Source().(*person); ok {
						return p.id, nil
					}
					return nil, nil
				}),
			gql.
				Field("data").
				Type(gql.List(dd)).
				Resolve(func(r interfaces.Resolver) (interface{}, error) {
					return dataArray, nil
				}),
		)

	fi := gql.
		Field("friend").
		Type(ur).
		Resolve(func(r interfaces.Resolver) (interface{}, error) {
			if p, ok := r.Source().(*person); ok {
				return &person{id: p.id + 1}, nil
			}
			return nil, nil
		})

	ur.RecursiveFields(fi)

	sc := gql.Schema(
		gql.Subscription(
			gql.
				Root().
				Fields(
					gql.
						Field("note").
						Type(nt).
						Args(
							gql.
								Arg("id").
								Type(gql.ID).
								Require(true),
						).
						Resolve(func(r interfaces.Resolver) (interface{}, error) {
							if b, ok := r.Source().([]byte); ok {
								var note = &note{}
								json.Unmarshal(b, note)
								return note, nil
							}
							return nil, nil
						}),
				),
		),
		gql.Query(
			gql.
				Root().
				Fields(
					gql.
						Field("note").
						Type(nt).
						Args(
							gql.
								Arg("id").
								Type(gql.ID).
								Require(true),
						).
						Resolve(func(r interfaces.Resolver) (interface{}, error) {
							for _, note := range notes {
								if note.ID == r.MustArg("id").String() {
									return note, nil
								}
							}
							return nil, nil
						}),
					gql.
						Field("user").
						Type(ur).
						Resolve(func(r interfaces.Resolver) (interface{}, error) {
							return &person{id: 0}, nil
						}),
					gql.
						Field("hello").
						Type(gql.String).
						Args(
							gql.
								Arg("input").
								Type(gql.String).
								Default("a default input").
								Require(false),
							gql.
								Arg("test").
								Require(true).
								Type(
									gql.
										Object("HelloInput").
										Args(
											gql.
												Arg("messageA").
												Description("the message A").
												Type(gql.String).
												Require(false),
											gql.
												Arg("messageB").
												Description("the message B").
												Type(gql.String).
												Default("hello world").
												Require(false),
										),
								),
						).
						Resolve(func(r interfaces.Resolver) (interface{}, error) {
							return r.MustArg("test").In("messageB").String(), nil
						}),
				),
		),
		gql.Mutation(
			gql.
				Root().
				Fields(
					gql.
						Field("addNote").
						Type(nt).
						Args(
							gql.
								Arg("content").
								Type(gql.String).
								Require(true),
						).
						Resolve(func(r interfaces.Resolver) (interface{}, error) {
							id := strconv.Itoa(int(time.Now().Unix()))
							et := &note{
								ID:      id,
								Content: r.MustArg("content").String(),
							}
							notes = append(notes, et)
							return et, nil
						}),
					gql.
						Field("updateNote").
						Type(nt).
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
							for _, note := range notes {
								if note.ID == r.MustArg("id").String() {
									note.Content = r.MustArg("content").String()
									o, _ := json.Marshal(note)
									r.Context().
										GQLSubscription().
										Emit(
											event.New(
												event.Field("note"),
												event.Payload(o),
												event.Filters(map[string]interface{}{
													"id": note.ID,
												}),
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

	gw.Schema(sc)

	ro.
		Post("/graphql").
		Params(
			append(GQLQueries, GQLBodies...)...,
		).
		Handle(GraphQL(sc))

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
