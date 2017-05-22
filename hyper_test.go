package hyper

import (
	"testing"

	"github.com/samuelngs/hyper/router"
	"github.com/samuelngs/hyper/sync"
)

func TestNew(t *testing.T) {

	h := New(
		Addr(":4000"),
		HTTP2(),
	)

	ws := h.Sync()

	ws.Namespace("default").
		Alias("test").
		Name("DefaultNamespace").
		Doc(`Default websocket namespace`).
		Summary(`Default websocket namespace`).
		Authorize(func(c sync.Context) error {
			return nil
		}).
		Middleware(func(d []byte, c sync.Context) {
		}).
		Handle("ping", func(d []byte, c sync.Context) {
		}).
		Handle("pong", func(d []byte, c sync.Context) {
		}).
		Catch(func(d []byte, c sync.Context) {
		})

	ro := h.Router()

	ro.Get("/").
		Alias("/test").
		Name("TestIndex").
		Doc(`Test index page`).
		Summary(`Test index page`).
		Params(
			Query("greeting").
				Doc(`The greeting message`).
				Summary(`The greeting message`).
				Default([]byte("Hello")).
				Require(false),
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
			c.Write([]byte("!"))
		}).
		Catch(func(c router.Context) {
			c.Error(c.Recover())
		})

	h.Run()

}
