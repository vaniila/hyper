package pubsub

// This code is for educational purposes only. It demonstrates a
// basic implemention of custom pubsub message module. It is not secure.
// Do not use it in production.

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/vaniila/hyper/message"
)

type body struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

type handler struct {
	name string
	fn   message.Handler
}

type server struct {
	client   *redis.Client
	pubsub   *redis.PubSub
	handlers []*handler
}

func (v *server) Start() error {
	v.client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	v.pubsub = v.client.Subscribe("default")
	go func() {
		for {
			msg, err := v.pubsub.ReceiveMessage()
			if err != nil {
				return
			}
			var b = new(body)
			if err := json.Unmarshal([]byte(msg.Payload), &b); err == nil {
				message := []byte(b.Message)
				for _, handler := range v.handlers {
					if handler.name == b.Channel {
						handler.fn(message)
					}
				}
			}
		}
	}()
	return nil
}

func (v *server) Stop() error {
	if err := v.pubsub.Close(); err != nil {
		return err
	}
	return v.client.Close()
}

func (v *server) Emit(channel []byte, message []byte) error {
	b := &body{
		Channel: string(channel),
		Message: string(message),
	}
	j, err := json.Marshal(b)
	if err != nil {
		return err
	}
	if _, err := v.client.Publish("default", string(j)).Result(); err != nil {
		return err
	}
	return nil
}

func (v *server) Listen(channel []byte, fn message.Handler) message.Close {
	hr := &handler{
		name: string(channel),
		fn:   fn,
	}
	v.handlers = append(v.handlers, hr)
	return func() {
		for i, h := range v.handlers {
			if hr == h {
				v.handlers = append(v.handlers[:i], v.handlers[i+1:]...)
			}
		}
	}
}

func (v *server) String() string {
	return "Redis::PubSub"
}

// New creates pubsub
func New() message.Service {
	s := &server{
		handlers: make([]*handler, 0),
	}
	return s
}
