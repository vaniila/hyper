package field

import (
	"net/http"
	"time"

	"github.com/vaniila/hyper/fault"
	"github.com/vaniila/hyper/gql/interfaces"
	"github.com/vaniila/hyper/router"
)

// value struct
type value struct {
	fmt    int
	key    string
	val    []byte
	has    bool
	parsed interface{}
}

func (v *value) In(s string) interfaces.Value {
	switch o := v.parsed.(type) {
	case []interfaces.Value:
		for _, r := range o {
			if r.Key() == s {
				return r
			}
		}
	}
	return nil
}

// Key of value
func (v *value) Key() string {
	return v.key
}

// Val is the data of value
func (v *value) Val() []byte {
	return v.val
}

// ThrowIllegal throws illegal type error
func (v *value) ThrowIllegal(cause fault.Cause) {
	err := fault.
		New("Illegal Action").
		SetStatus(http.StatusInternalServerError).
		AddCause(cause)
	panic(err)
}

// MustInt returns the data in int format
func (v *value) MustInt() int {
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
func (v *value) MustI32() int32 {
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
func (v *value) MustI64() int64 {
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
func (v *value) MustU32() uint32 {
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
func (v *value) MustU64() uint64 {
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
func (v *value) MustF32() float32 {
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
func (v *value) MustF64() float64 {
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
func (v *value) MustBool() bool {
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
func (v *value) MustTime() time.Time {
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

// MustArray returns value in array format
func (v *value) MustArray() []interface{} {
	return v.parsed.([]interface{})
}

// Any returns value
func (v *value) Any() interface{} {
	return v.parsed
}

// Has represents if input exists
func (v *value) Has() bool {
	return v.has
}

func (v *value) String() string {
	return string(v.val)
}
