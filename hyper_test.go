package hyper

import (
	"log"
	"testing"

	"github.com/samuelngs/hyper/router"
)

type TestHTTPResponse struct {
	Message string
}

func TestNew(t *testing.T) {

	h := New(
		Addr(":4000"),
		HTTP2(),
	)

	ro := h.Router()

	ro.Get("/").
		Alias("/:name").
		Name("GetUsername").
		Doc(`Retrieve username`).
		Summary(`Retrieve username`).
		Params(
			Query("name").
				Doc(`The username`).
				Summary(`The username`).
				Default([]byte("")).
				Require(true),
		).
		Models(
			Model(StatusOK, new(string)),
			Model(StatusOK, new(TestHTTPResponse)),
			Model(StatusForbidden, new(TestHTTPResponse)),
		).
		Middleware(func(c router.Context) {
			c.Write([]byte("uid => "))
			c.Write(c.Cookie().MustGet("name").Val())
		}).
		Handle(func(c router.Context) {
			c.Write([]byte(" | "))
			c.Write([]byte(c.ProcessID()))
		}).
		Catch(func(c router.Context) {
			c.Error(c.Recover())
		})

	ns := ro.Namespace("/test").
		Name("testing").
		Doc(`testing documentation`).
		Summary(`testing summary`)

	ns.Get("/").
		Name("testing").
		Doc(``).
		Summary(``).
		Handle(func(c router.Context) {
			c.Status(200).Write([]byte("hello world"))
		})

	no := ns.Namespace("/hello").
		Name("hello").
		Doc(`hello documentation`).
		Summary(`hello summary`)

	no.Get("/").
		Name("hello").
		Doc(``).
		Summary(``).
		Handle(func(c router.Context) {
			c.Status(200).Write([]byte("hello"))
		})

	log.Print(ro)

	if e := h.Run(); e != nil {
		log.Print(e)
	}

}
