package enum

type valueconfig struct {
	value *value
}

func (v *valueconfig) Name() string {
	return v.value.name
}

func (v *valueconfig) Description() string {
	return v.value.description
}

func (v *valueconfig) Is() interface{} {
	return v.value.value
}

func (v *valueconfig) Deprecation() string {
	return v.value.deprecationReason
}
