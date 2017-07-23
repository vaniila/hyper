package engine

type identity struct {
	hasid, haskey, hasprofile, hasmachine bool
	id                                    int
	key, profile, machine                 string
}

func (v *identity) HasID() bool {
	return v.hasid
}

func (v *identity) GetID() int {
	return v.id
}

func (v *identity) SetID(i int) {
	v.id = i
	v.hasid = true
}

func (v *identity) HasKey() bool {
	return v.haskey
}

func (v *identity) GetKey() string {
	return v.key
}

func (v *identity) SetKey(s string) {
	v.key = s
	v.haskey = true
}

func (v *identity) HasProfile() bool {
	return v.hasprofile
}

func (v *identity) GetProfile() string {
	return v.profile
}

func (v *identity) SetProfile(s string) {
	v.profile = s
	v.hasprofile = true
}

func (v *identity) HasMachine() bool {
	return v.hasmachine
}

func (v *identity) GetMachine() string {
	return v.machine
}

func (v *identity) SetMachine(s string) {
	v.machine = s
	v.hasmachine = true
}
