package hyper

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/julienbreux/hyper/logger" // FIXME: Before proposal
	"github.com/vaniila/hyper/cache"
	"github.com/vaniila/hyper/dataloader"
	"github.com/vaniila/hyper/engine"
	"github.com/vaniila/hyper/gws"
	"github.com/vaniila/hyper/message"
	"github.com/vaniila/hyper/router"
	"github.com/vaniila/hyper/sync"
	"github.com/vaniila/hyper/websocket"
)

// Hyper service
type Hyper struct {
	id         string
	addr       string
	logger     logger.Service
	cache      cache.Service
	message    message.Service
	dataloader dataloader.Service
	router     router.Service
	engine     engine.Service
	sync       sync.Service
	gws        gws.Service
	websocket  websocket.Service
}

func (v *Hyper) start() error {
	if err := v.logger.Start(); err != nil {
		return err
	}
	if err := v.cache.Start(); err != nil {
		return err
	}
	if err := v.message.Start(); err != nil {
		return err
	}
	if err := v.dataloader.Start(); err != nil {
		return err
	}
	if err := v.router.Start(); err != nil {
		return err
	}
	if err := v.sync.Start(); err != nil {
		return err
	}
	if err := v.gws.Start(); err != nil {
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
	if err := v.gws.Stop(); err != nil {
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
	if err := v.dataloader.Stop(); err != nil {
		return err
	}
	if err := v.logger.Stop(); err != nil {
		return err
	}
	return nil
}

// Run hyper server
func (v *Hyper) Run() error {

	if err := v.start(); err != nil {
		return err
	}

	v.logger.Info("Server running on %s", logger.NewField("address", v.addr))

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	v.logger.Info("Received signal %s", logger.NewField("signal", <-ch))

	return v.stop()
}

// Logger returns logger service
func (v *Hyper) Logger() logger.Service {
	return v.logger
}

// Sync returns sync service
func (v *Hyper) Sync() sync.Service {
	return v.sync
}

// Gws returns graphql subscription service
func (v *Hyper) Gws() gws.Service {
	return v.gws
}

// Router returns router service
func (v *Hyper) Router() router.Service {
	return v.router
}

// String returns hyper server in string
func (v *Hyper) String() string {
	return "Hyper::Server"
}
