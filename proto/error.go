package proto

import (
	"fmt"
)

// Error struct
type Error struct {
	Code   uint16
	Class  uint16
	Method uint16
	Msg    string
	Soft   bool
}

// NewSoftError returns a soft error
func NewSoftError(code uint16, msg string, class uint16, method uint16) *Error {
	return &Error{
		Code:   code,
		Class:  class,
		Method: method,
		Msg:    msg,
		Soft:   true,
	}
}

// NewHardError returns an fatal error
func NewHardError(code uint16, msg string, class uint16, method uint16) *Error {
	return &Error{
		Code:   code,
		Class:  class,
		Method: method,
		Msg:    msg,
		Soft:   false,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("Exception (%d) Reason: %q - for class (%d) method (%d)", e.Code, e.Msg, e.Class, e.Method)
}
