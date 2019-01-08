package fault

import "fmt"

// Template struct
type Formatter string

func (v Formatter) Fill(a ...interface{}) Context {
	s := fmt.Sprintf(string(v), a...)
	return New(s)
}
