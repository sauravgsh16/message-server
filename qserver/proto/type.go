package proto

import (
        "io"
)

type Frame interface {
        FrameType() byte
}

type MethodFrame interface {
        MethodName() string
        MethodIdentifier() (uint16, uint16)
        Read(r io.Reader) (err error)
        Write(w io.Writer) (err error)
        FrameType() byte
}