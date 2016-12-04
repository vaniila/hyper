package hyper

import "github.com/samuelngs/hyper/router"

type param struct {
	typ           router.ParamType
	format        router.DataFormat
	name          string
	summary       string
	documentation string
	defaults      []byte
	require       bool
}

func (v *param) Format(f router.DataFormat) router.Param {
	v.format = f
	return v
}

func (v *param) Summary(s string) router.Param {
	v.summary = s
	return v
}

func (v *param) Doc(s string) router.Param {
	v.documentation = s
	return v
}

func (v *param) Default(b []byte) router.Param {
	v.defaults = b
	return v
}

func (v *param) Require(b bool) router.Param {
	v.require = b
	return v
}

// Query func
func Query(name string) router.Param {
	return &param{}
}

// Body func
func Body(name string) router.Param {
	return &param{}
}

// Param func
func Param(name string) router.Param {
	return &param{}
}

// Header func
func Header(name string) router.Param {
	return &param{}
}
