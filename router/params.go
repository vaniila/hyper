package router

// ParamType type
type ParamType int

// DataFormat type
type DataFormat int

// UnknownType type
const (
	UnknownParam ParamType = iota
	ParamBody
	ParamParam
	ParamQuery
	ParamHeader
	ParamCookie

	UnknownData DataFormat = iota
	String
	Number
	Float
	Bool
	Binary
	Email
	Phone
	Latitude
	Longitude
	Date
	Time
	UUID
	JSON
)
