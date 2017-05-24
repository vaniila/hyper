package fault

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

// Context interface
type Context interface {
	ID() string
	Status() int
	Message() string
	Causes() []Cause
	Error() string
	SetStatus(int) Context
	AddCause(...Cause) Context
	Json() []byte
	JsonString() string
}

// Cause interface
type Cause interface {
	Resource() string
	Field() string
	Code() string
	SetResource(string) Cause
	SetField(string) Cause
	SetCode(string) Cause
}

// Template interface
type Template interface {
	Fill(...interface{}) Context
}

func newID() string {
	b := new([16]byte)
	rand.Read(b[:])
	b[8] = (b[8] | 0x40) & 0x7F
	b[6] = (b[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// New creates fault error from string
func New(message string) Context {
	return &Fault{
		FID:      newID(),
		FStatus:  http.StatusForbidden,
		FMessage: message,
		FCauses:  make([]Cause, 0),
	}
}

// Wrap creates fault error from error
func Wrap(err error) Context {
	if err != nil {
		return New(err.Error())
	}
	return nil
}

// Format to compose failure reusable template
func Format(format string) Template {
	return Formatter(format)
}

// For to create error
func For(code string) Cause {
	return &Reason{
		RCode: code,
	}
}

// Is to check if error is fault context
func Is(e error) (Context, bool) {
	v, ok := e.(Context)
	return v, ok
}

// Match error
func Match() {
}
