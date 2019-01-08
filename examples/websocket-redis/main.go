package main

// This code is for educational purposes only. It demonstrates a
// basic implemention of custom pubsub message module. It is not secure.
// Do not use it in production.

import (
	"github.com/vaniila/hyper"
	"github.com/vaniila/hyper/examples/websocket-redis/pubsub"
	"github.com/vaniila/hyper/sync"
)

func main() {

	redisPubsub := pubsub.New()

	h := hyper.New(
		hyper.Addr(":4000"),
		hyper.HTTP2(),
		hyper.Message(redisPubsub),
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

	h.Run()
}
