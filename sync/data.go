package sync

type data struct {
	typ   int32
	bytes []byte
}

func (v *data) Type() int32 {
	return v.typ
}

func (v *data) Bytes() []byte {
	return v.bytes
}
