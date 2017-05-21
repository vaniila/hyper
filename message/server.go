package message

type event struct {
	channel []byte
	message []byte
}

type server struct {
	id     string
	events chan *event
}

func (v *server) Start() error {
	v.events = make(chan *event)
	return nil
}

func (v *server) Stop() error {
	close(v.events)
	return nil
}

func (v *server) Emit(channel, message []byte) error {
	e := &event{
		channel: channel,
		message: message,
	}
	v.events <- e
	return nil
}

func (v *server) Listen(channel []byte) (<-chan []byte, chan<- struct{}, error) {
	c := make(chan []byte)
	s := make(chan struct{})
	go func() {
		defer func() {
			close(c)
			close(s)
		}()
		for {
			select {
			case b := <-v.events:
				if b != nil && string(b.channel[:]) == string(channel[:]) {
					c <- b.message
				}
			case <-s:
				break
			}
		}
	}()
	return c, s, nil
}

func (v *server) String() string {
	return "Hyper::Message"
}
