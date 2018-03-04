package main

import (
	"fmt"
	"net/http"

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
		Post("/").
		Params(
			hyper.Body("file").
				Format(hyper.File).
				Require(true),
		).
		Handle(func(c router.Context) {
			file := c.MustBody("file").Val()
			typ := http.DetectContentType(file)
			msg := fmt.Sprintf("File uploaded: %s type, %d bytes", typ, len(file))
			c.Write([]byte(msg))
		})

	h.Run()
}
