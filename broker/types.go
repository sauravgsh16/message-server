package broker

import (
        "io"
        "fmt"
)

const (
        FrameError   = 501
        ChannelError = 504
)

var (
        ErrFrame = &Error{Code: FrameError, Reason: "frame could not be parsed"}
        ErrClosed = &Error{Code: ChannelError, Reason: "channel/connection is not open"}
)

type message interface {
      read(io.Reader) error
      write(io.Writer) error  
}

type frame interface {
        write(io.Writer) error
}

type writer struct {
        w io.Writer
}

type reader struct {
        r io.Reader
}

type exchangeDeclare struct {
        Exchange string
        Type     string
}

type exchangeDeclareOk struct {}


type bodyFrame struct {
        ChannelId uint16
        Body      []byte
}

func (f *bodyFrame) channel() uint16 { return f.ChannelId }


type Error struct {
        Code   int     // constant code
        Reason string  // description of the error
}

func (e Error) Error() string {
        return fmt.Sprintf("Exception (%d) Reason: %q", e.Code, e.Reason)
} 