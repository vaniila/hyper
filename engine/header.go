package engine

import "net/http"

// Header struct
type Header struct {
	header http.Header
}

// Set data to response header
func (v *Header) Set(key string, val string) {
}

// Get data from request header
func (v *Header) Get(key string) *Value {
	return nil
}
