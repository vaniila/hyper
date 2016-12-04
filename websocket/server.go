package websocket

type server struct {
	id string
}

func (v *server) Start() error {
	return nil
}

func (v *server) Stop() error {
	return nil
}

func (v *server) String() string {
	return "Hyper::Cache"
}
