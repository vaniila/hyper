package field

import (
	"net/http"
	"time"

	"github.com/vaniila/hyper/fault"
	"github.com/vaniila/hyper/router"
)

// Value struct
type Value struct {
	fmt    int
	key    string
	val    []byte
	has    bool
	parsed interface{}
}

// Type of value
func (v *Value) Type() router.ParamType {
	return router.UnknownParam
}

// Key of value
func (v *Value) Key() string {
	return v.key
}

// Val is the data of value
func (v *Value) Val() []byte {
	return v.val
}

// ThrowIllegal throws illegal type error
func (v *Value) ThrowIllegal(cause fault.Cause) {
	err := fault.
		New("Illegal Action").
		SetStatus(http.StatusInternalServerError).
		AddCause(cause)
	panic(err)
}

// MustInt returns the data in int format
func (v *Value) MustInt() int {
	if v.fmt != router.Int {
		v.ThrowIllegal(
			fault.
				For(fault.Illegal).
				SetResource("MustInt").
				SetField(v.key),
		)
	}
	return v.parsed.(int)
}

// MustI32 returns data in int32 format
func (v *Value) MustI32() int32 {
	if v.fmt != router.I32 {
		v.ThrowIllegal(
			fault.
				For(fault.Illegal).
				SetResource("MustI32").
				SetField(v.key),
		)
	}
	return v.parsed.(int32)
}

// MustI64 returns data in int64 format
func (v *Value) MustI64() int64 {
	if v.fmt != router.I64 {
		v.ThrowIllegal(
			fault.
				For(fault.Illegal).
				SetResource("MustI64").
				SetField(v.key),
		)
	}
	return v.parsed.(int64)
}

// MustU32 returns data in uint32 format
func (v *Value) MustU32() uint32 {
	if v.fmt != router.U32 {
		v.ThrowIllegal(
			fault.
				For(fault.Illegal).
				SetResource("MustU32").
				SetField(v.key),
		)
	}
	return v.parsed.(uint32)
}

// MustU64 returns data in uint64 format
func (v *Value) MustU64() uint64 {
	if v.fmt != router.U64 {
		v.ThrowIllegal(
			fault.
				For(fault.Illegal).
				SetResource("MustU64").
				SetField(v.key),
		)
	}
	return v.parsed.(uint64)
}

// MustF32 returns data in float32 format
func (v *Value) MustF32() float32 {
	if v.fmt != router.F32 {
		v.ThrowIllegal(
			fault.
				For(fault.Illegal).
				SetResource("MustF32").
				SetField(v.key),
		)
	}
	return v.parsed.(float32)
}

// MustF64 returns data in float64 format
func (v *Value) MustF64() float64 {
	switch v.fmt {
	case router.F64, router.Lat, router.Lon:
		// whitelist
	default:
		v.ThrowIllegal(
			fault.
				For(fault.Illegal).
				SetResource("MustF64").
				SetField(v.key),
		)
	}
	return v.parsed.(float64)
}

// MustBool returns data in boolean format
func (v *Value) MustBool() bool {
	if v.fmt != router.Bool {
		v.ThrowIllegal(
			fault.
				For(fault.Illegal).
				SetResource("MustBool").
				SetField(v.key),
		)
	}
	return v.parsed.(bool)
}

// MustTime returns data in time.Time format
func (v *Value) MustTime() time.Time {
	switch v.fmt {
	case router.DateTimeRFC822, router.DateTimeRFC3339, router.DateTimeUnix:
		// whitelist
	default:
		v.ThrowIllegal(
			fault.
				For(fault.Illegal).
				SetResource("MustTime").
				SetField(v.key),
		)
	}
	return v.parsed.(time.Time)
}

// Has represents if input exists
func (v *Value) Has() bool {
	return v.has
}

func (v *Value) String() string {
	return string(v.val)
}
