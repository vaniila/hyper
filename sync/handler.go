package sync

type handler struct {
	name string
	fn   HandlerFunc
}

func (v *handler) Name() string {
	return v.name
}

func (v *handler) Func() HandlerFunc {
	return v.fn
}
