package broker

import (
        "io"
        "fmt"
)

const (
        frameEnd     = 206
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

type messageWithContent interface {
        message
        getContent() (properties, []byte)
        setContent(properties, []byte)
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

type methodFrame struct {
        ChannelId uint16
        Method    message
}

type headerFrame struct {
        ChannelId  uint16
        Size       uint64
        Properties properties
}

type bodyFrame struct {
        ChannelId uint16
        Body      []byte
}

func (f *bodyFrame) channel() uint16 { return f.ChannelId }

type Publishing struct {
        UserId    string
        MessageId string
        Body      []byte
}

type properties struct {
        UserId    string
        MessageId string
}

type Error struct {
        Code   int     // constant code
        Reason string  // description of the error
}

func (e Error) Error() string {
        return fmt.Sprintf("Exception (%d) Reason: %q", e.Code, e.Reason)
} 