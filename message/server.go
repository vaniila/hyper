package message

import "bytes"

type handler struct {
	name []byte
	fn   Handler
}

type server struct {
	id       string
	handlers []*handler
}

func (v *server) Start() error {
	return nil
}

func (v *server) Stop() error {
	return nil
}

func (v *server) Emit(channel, message []byte) error {
	for _, handle := range v.handlers {
		if bytes.Equal(handle.name, channel) {
			handle.fn(message)
		}
	}
	return nil
}

func (v *server) Listen(channel []byte, fn Handler) Close {
	hr := &handler{
		name: channel,
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
	return "Hyper::Message"
}
