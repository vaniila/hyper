package main

import (
	"strings"

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
			hyper.OneOf(
				hyper.Query("prefix").
					Format(hyper.Text).
					Default([]byte("suffix-")).
					Require(false),
				hyper.Query("suffix").
					Format(hyper.Text).
					Default([]byte("-prefix")).
					Require(false),
			),
		).
		Handle(func(c router.Context) {
			var parts []string
			if p := c.MustQuery("prefix"); p.Has() {
				parts = append(parts, p.String())
			}
			parts = append(parts, c.MustQuery("message").String())
			if p := c.MustQuery("suffix"); p.Has() {
				parts = append(parts, p.String())
			}
			c.Write([]byte(strings.Join(parts, "")))
		})

	h.Run()
}
