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
		Alias("/oh").
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
		}).
		Handle(func(c router.Context) {
			c.Write([]byte(" | "))
			c.Write([]byte(c.ProcessID()))
		}).
		Catch(func(c router.Context) {
			c.Error(c.Recover())
		})

	log.Print(ro)

	if e := h.Run(); e != nil {
		log.Print(e)
	}

}
