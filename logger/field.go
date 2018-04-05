package logger

type field struct {
	key   string
	value interface{}
}

// Value returns the key of field
func (f *field) Key() string {
	return f.key
}

// Value returns the value of field
func (f *field) Value() interface{} {
	return f.value
}
