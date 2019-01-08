package dataloader

// Return type
type Return struct {
	data interface{}
	err  error
}

// Returns type
type Returns struct {
	data []interface{}
	errs []error
}

// Data returns data of return object
func (v *Return) Data() interface{} {
	return v.data
}

// Error returns the error of return object
func (v *Return) Error() error {
	return v.err
}

// Data returns items of data within returns object
func (v *Returns) Data() []interface{} {
	return v.data
}

// Errors returns array of errors within returns object
func (v *Returns) Errors() []error {
	return v.errs
}
