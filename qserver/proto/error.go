package proto

type ProtoError struct {
        Code   uint16
        Class  uint16
        Method uint16
        Msg    string
        Soft   bool
}

func NewSoftError(code uint16, msg string, class uint16, method uint16) *ProtoError {
        return &ProtoError{
                Code:   code,
                Class:  class,
                Method: method,
                Msg:    msg,
                Soft:   true,
        }
}


func NewHardError(code uint16, msg string, class uint16, method uint16) *ProtoError {
        return &ProtoError{
                Code:   code,
                Class:  class,
                Method: method,
                Msg:    msg,
                Soft:   false,
        }
}