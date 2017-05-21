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

type paramconfig struct {
	*param
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

func (v *param) Config() router.ParamConfig {
	return &paramconfig{v}
}

func (v *paramconfig) Name() string {
	return v.name
}

func (v *paramconfig) Type() router.ParamType {
	return v.typ
}

func (v *paramconfig) Format() router.DataFormat {
	return v.format
}

func (v *paramconfig) Summary() string {
	return v.summary
}

func (v *paramconfig) Doc() string {
	return v.documentation
}

func (v *paramconfig) Default() []byte {
	return v.defaults
}

func (v *paramconfig) Require() bool {
	return v.require
}

// Query func
func Query(name string) router.Param {
	return &param{typ: router.ParamQuery, name: name}
}

// Body func
func Body(name string) router.Param {
	return &param{typ: router.ParamBody, name: name}
}

// Param func
func Param(name string) router.Param {
	return &param{typ: router.ParamParam, name: name}
}

// Header func
func Header(name string) router.Param {
	return &param{typ: router.ParamHeader, name: name}
}

// Cookie func
func Cookie(name string) router.Param {
	return &param{typ: router.ParamCookie, name: name}
}
