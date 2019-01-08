package router

import (
	"encoding/json"
	"math"
	"regexp"
	"strconv"
	"time"
)

const (
	Any int = iota
	Text
	Int
	I32
	I64
	U32
	U64
	F32
	F64
	Bool
	Lat
	Lon
	Json
	Email
	URL
	UUID
	File
	DateTimeRFC822
	DateTimeRFC3339
	DateTimeUnix
)

var validations = map[int]func([]byte) (interface{}, bool){
	Any:             ok,
	Text:            ok,
	Int:             isInt,
	I32:             isI32,
	I64:             isI64,
	U32:             isU32,
	U64:             isU64,
	F32:             isF32,
	F64:             isF64,
	Bool:            isBool,
	Lat:             isLat,
	Lon:             isLon,
	Json:            isJson,
	Email:           isEmail,
	URL:             isURL,
	UUID:            isUUID,
	File:            ok,
	DateTimeRFC822:  isDateTimeRFC822,
	DateTimeRFC3339: isDateTimeRFC3339,
	DateTimeUnix:    isDateTimeUnix,
}

func Val(typ int, v []byte) (interface{}, bool) {
	validation, ok := validations[typ]
	if !ok {
		return v, true
	}
	return validation(v)
}

func ok(v []byte) (interface{}, bool) {
	return v, true
}

func isInt(v []byte) (interface{}, bool) {
	s := string(v[:])
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, false
	}
	if i < math.MinInt64 || i > math.MaxInt64 {
		return nil, false
	}
	return i, true
}

func isI32(v []byte) (interface{}, bool) {
	s := string(v[:])
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, false
	}
	if i < math.MinInt32 || i > math.MaxInt32 {
		return nil, false
	}
	return int32(i), true
}

func isI64(v []byte) (interface{}, bool) {
	s := string(v[:])
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, false
	}
	if i < math.MinInt64 || i > math.MaxInt64 {
		return nil, false
	}
	return int64(i), true
}

func isU32(v []byte) (interface{}, bool) {
	s := string(v[:])
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, false
	}
	if i < 0 || i > math.MaxUint32 {
		return nil, false
	}
	return uint32(i), true
}

func isU64(v []byte) (interface{}, bool) {
	s := string(v[:])
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, false
	}
	if i < 0 || i > math.MaxUint64 {
		return nil, false
	}
	return i, true
}

func isF32(v []byte) (interface{}, bool) {
	s := string(v[:])
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, false
	}
	if f < -math.MaxFloat32 || f > math.MaxFloat32 {
		return nil, false
	}
	return float32(f), true
}

func isF64(v []byte) (interface{}, bool) {
	s := string(v[:])
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, false
	}
	if f < -math.MaxFloat64 || f > math.MaxFloat64 {
		return nil, false
	}
	return f, true
}

func isBool(v []byte) (interface{}, bool) {
	s := string(v[:])
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil, false
	}
	return b, true
}

func isLat(v []byte) (interface{}, bool) {
	r := regexp.
		MustCompile("^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$").
		MatchString(string(v[:]))
	if !r {
		return nil, false
	}
	f, ok := isF64(v)
	if !ok {
		return nil, false
	}
	return f.(float64), true
}

func isLon(v []byte) (interface{}, bool) {
	r := regexp.
		MustCompile("^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$").
		MatchString(string(v[:]))
	if !r {
		return nil, false
	}
	f, ok := isF64(v)
	if !ok {
		return nil, false
	}
	return f.(float64), true
}

func isJson(v []byte) (interface{}, bool) {
	var tmp json.RawMessage
	if json.Unmarshal(v, &tmp) != nil {
		return nil, false
	}
	return tmp, true
}

func isEmail(v []byte) (interface{}, bool) {
	r := regexp.
		MustCompile("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:\\(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22)))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$").
		MatchString(string(v[:]))
	if !r {
		return nil, false
	}
	return v, true
}

func isURL(v []byte) (interface{}, bool) {
	r := regexp.
		MustCompile("((([A-Za-z]{3,9}:(?://)?)(?:[-;:&=+$,\\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=+$,\\w]+@)[A-Za-z0-9.-]+)((?:/[+~%/.\\w-_]*)?\\??(?:[-+=&;%@.\\w_]*)#?(?:[\\w]*))?)").
		MatchString(string(v[:]))
	if !r {
		return nil, false
	}
	return v, true
}

func isUUID(v []byte) (interface{}, bool) {
	r := regexp.
		MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$").
		MatchString(string(v[:]))
	if !r {
		return nil, false
	}
	return v, true
}

func isDateTimeRFC822(v []byte) (interface{}, bool) {
	s := string(v[:])
	if t, err := time.Parse(time.RFC822Z, s); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.RFC822, s); err == nil {
		return t, true
	}
	return nil, false
}

func isDateTimeRFC3339(v []byte) (interface{}, bool) {
	s := string(v[:])
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, true
	}
	return nil, false
}

func isDateTimeUnix(v []byte) (interface{}, bool) {
	s := string(v[:])
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, false
	}
	t := time.Unix(0, i*int64(time.Millisecond))
	if t.IsZero() {
		return nil, false
	}
	return t, true
}
