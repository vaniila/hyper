package hyper

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/samuelngs/hyper/cache"
	"github.com/samuelngs/hyper/engine"
	"github.com/samuelngs/hyper/message"
	"github.com/samuelngs/hyper/router"
	"github.com/samuelngs/hyper/sync"
	"github.com/samuelngs/hyper/websocket"
)

// Hyper service
type Hyper struct {
	id        string
	addr      string
	cache     cache.Service
	message   message.Service
	router    router.Service
	engine    engine.Service
	sync      sync.Service
	websocket websocket.Service
}

func (v *Hyper) start() error {
	if err := v.cache.Start(); err != nil {
		return err
	}
	if err := v.message.Start(); err != nil {
		return err
	}
	if err := v.router.Start(); err != nil {
		return err
	}
	if err := v.sync.Start(); err != nil {
		return err
	}
	if err := v.websocket.Start(); err != nil {
		return err
	}
	if err := v.engine.Start(); err != nil {
		return err
	}
	return nil
}

// Stop hyper server
func (v *Hyper) stop() error {
	if err := v.engine.Stop(); err != nil {
		return err
	}
	if err := v.websocket.Stop(); err != nil {
		return err
	}
	if err := v.sync.Stop(); err != nil {
		return err
	}
	if err := v.router.Stop(); err != nil {
		return err
	}
	if err := v.cache.Stop(); err != nil {
		return err
	}
	if err := v.message.Stop(); err != nil {
		return err
	}
	return nil
}

// Run hyper server
func (v *Hyper) Run() error {

	if err := v.start(); err != nil {
		return err
	}

	log.Printf("Server running on %s", v.addr)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	log.Printf("Received signal %s", <-ch)

	return v.stop()
}

// Sync returns sync service
func (v *Hyper) Sync() sync.Service {
	return v.sync
}

// Router returns router service
func (v *Hyper) Router() router.Service {
	return v.router
}

// String returns hyper server in string
func (v *Hyper) String() string {
	return "Hyper::Server"
}
