package sync

type identity struct {
	hasid  bool
	haskey bool
	id     int
	key    string
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
