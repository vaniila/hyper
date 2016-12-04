package engine

import "net/http"

// Cookie struct
type Cookie struct {
	req *http.Request
	res http.ResponseWriter
}

// Set data to response cookie
func (v *Cookie) Set(key string, val string) {
}

// Get data from request cookie
func (v *Cookie) Get(key string) *Value {
	return nil
}
