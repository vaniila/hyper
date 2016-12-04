package tracer

// Tracer interface
type Tracer interface {
	StartSpan(name string) Span
	String() string
}

// Span interface
type Span interface {
	SetTag(key string, val string)
	SetOperationName(name string)
	SetBaggageItem(key, value string) Span
	BaggageItem(key string) string
	Tracer() Tracer
	LogEvent(e string)
	LogEventf(f string, v ...interface{})
	LogEventWithPayload(e string, payload interface{})
	LogKV(kv ...interface{})
	LogFields(f ...Field)
	Finish()
}

// Field interface
type Field interface {
	Key() string
	Value() interface{}
	String() string
}
