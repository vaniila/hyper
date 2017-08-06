package hyper

import (
	"log"

	"github.com/vaniila/hyper/router"
)

type param struct {
	typ           router.ParamType
	custom        router.CustomFunc
	format        int
	name          string
	summary       string
	documentation string
	defaults      []byte
	oneof         []router.Param
	require       bool
}

type paramconfig struct {
	*param
}

func (v *param) Custom(c router.CustomFunc) router.Param {
	v.custom = c
	return v
}

func (v *param) Format(f int) router.Param {
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

func (v *paramconfig) Custom() router.CustomFunc {
	return v.custom
}

func (v *paramconfig) Format() int {
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

func (v *paramconfig) OneOf() []router.Param {
	return v.oneof
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

// OneOf group func
func OneOf(ps ...router.Param) router.Param {
	for _, param := range ps {
		if param.Config().Require() {
			log.Fatalf("cannot set %v field as required, [oneof] parameters do not support required checking", param.Config().Name())
		}
	}
	return &param{typ: router.ParamOneOf, oneof: ps}
}
