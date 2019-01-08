package main

import (
	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/router"
)

func main() {

	h := hyper.New(
		hyper.Addr(":4000"),
		hyper.HTTP2(),
	)

	ro := h.Router()

	ro.
		Get("/").
		Params(
			hyper.Query("message").
				Format(hyper.Text).
				Doc(`custom greeting message`).
				Summary(`custom greeting message`).
				Default([]byte("hello world")).
				Require(false),
		).
		Handle(func(c router.Context) {
			c.Write(c.MustQuery("message").Val())
		})

	h.Run()
}
