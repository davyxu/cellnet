package cellnet

import "fmt"

type Error struct {
	s       string
	context interface{}
}

func (self *Error) Error() string {

	if self.context == nil {
		return self.s
	}

	return fmt.Sprintf("%s, '%v'", self.s, self.context)
}

func NewError(s string) error {

	return &Error{s: s}
}

func NewErrorContext(s string, context interface{}) error {
	return &Error{s: s, context: context}
}
