package sync

import (
	"bytes"
	"encoding/gob"
)

type packet struct {
	D int    // message direction, 0 => client to server, 1 => server to client
	T int    // message type
	N string // namespace
	C string // channel
	M []byte // bytes
}

func (v packet) Direction() int {
	return v.D
}

func (v packet) Type() int {
	return v.T
}

func (v packet) Namespace() string {
	return v.N
}

func (v packet) Channel() string {
	return v.C
}

func (v packet) Message() []byte {
	return v.M
}

func (v packet) Bytes() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decode(c []byte) (*packet, error) {
	var p *packet
	b := bytes.NewBuffer(c)
	dec := gob.NewDecoder(b)
	if err := dec.Decode(&p); err != nil {
		return nil, err
	}
	return p, nil
}
