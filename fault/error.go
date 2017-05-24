package fault

import "encoding/json"

// Fault struct
type Fault struct {
	FID      string  `json:"-"`
	FStatus  int     `json:"-"`
	FMessage string  `json:"message"`
	FCauses  []Cause `json:"causes"`
}

// ID of the fault
func (v *Fault) ID() string {
	return v.FID
}

// Status return status code
func (v *Fault) Status() int {
	return v.FStatus
}

// Message of the fault
func (v *Fault) Message() string {
	return v.FMessage
}

// Causes of the fault
func (v *Fault) Causes() []Cause {
	return v.FCauses
}

// SetStatus to set fault status
func (v *Fault) SetStatus(s int) Context {
	v.FStatus = s
	return v
}

// AddCause error
func (v *Fault) AddCause(s ...Cause) Context {
	for _, c := range s {
		v.FCauses = append(v.FCauses, c)
	}
	return v
}

// ErrorOrNil returns a fault error only if there are errors
func (v *Fault) ErrorOrNil() error {
	if len(v.FCauses) > 0 {
		return v
	}
	return nil
}

func (v *Fault) Error() string {
	return v.FMessage
}

// Json to return json bytes of the error
func (v *Fault) Json() []byte {
	b, _ := json.Marshal(v)
	return b
}

// JsonString to return json string of the error
func (v *Fault) JsonString() string {
	return string(v.Json()[:])
}
