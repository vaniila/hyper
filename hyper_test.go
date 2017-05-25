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
			n.Write(&sync.Packet{Message: []byte{49, 50, 51}})
		}).
		Catch(func(m []byte, n sync.Channel, c sync.Context) {
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
