package hyper

import "github.com/vaniila/hyper/router"

type model struct {
	code int
	hash string
}

func (v *model) Code() int {
	return v.code
}

func (v *model) Hash() string {
	return v.hash
}

// Model func
func Model(code int, response interface{}) router.Model {
	return &model{}
}
