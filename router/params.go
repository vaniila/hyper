package router

// ParamType type
type ParamType int

// CustomFunc type
type CustomFunc func(v []byte) bool

// UnknownType type
const (
	UnknownParam ParamType = iota
	ParamBody
	ParamParam
	ParamQuery
	ParamHeader
	ParamCookie
)

// String value for ParamType
func (v ParamType) String() string {
	switch v {
	case ParamBody:
		return "body"
	case ParamParam:
		return "param"
	case ParamQuery:
		return "query"
	case ParamHeader:
		return "header"
	case ParamCookie:
		return "cookie"
	}
	return "unknown"
}
