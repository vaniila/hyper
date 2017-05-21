package fault

// Some of the causes
var (
	Missing           = string("missing")
	MissingField      = string("missing_field")
	UnregisteredField = string("unregistered_field")
	Invalid           = string("invalid")
	AlreadyExists     = string("already_exists")
)

// Reason struct
type Reason struct {
	RResource string `json:"resource,omitempty"`
	RField    string `json:"field,omitempty"`
	RCode     string `json:"code,omitempty"`
}

// Resource to return resource name
func (v *Reason) Resource() string {
	return v.RResource
}

// Field to return field name
func (v *Reason) Field() string {
	return v.RField
}

// Code to return code name
func (v *Reason) Code() string {
	return v.RCode
}

// SetResource to set resource
func (v *Reason) SetResource(s string) Cause {
	v.RResource = s
	return v
}

// SetField to set field
func (v *Reason) SetField(s string) Cause {
	v.RField = s
	return v
}

// SetCode to set code
func (v *Reason) SetCode(s string) Cause {
	v.RCode = s
	return v
}
