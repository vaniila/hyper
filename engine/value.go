package engine

import "github.com/samuelngs/hyper/router"

// Value struct
type Value struct {
	typ router.ParamType
	key string
	val []byte
}

// Type of value
func (v *Value) Type() router.ParamType {
	return v.typ
}

// Key of value
func (v *Value) Key() string {
	return v.key
}

// Val is the data of value
func (v *Value) Val() []byte {
	return v.val
}

func (v *Value) String() string {
	return string(v.val)
}
