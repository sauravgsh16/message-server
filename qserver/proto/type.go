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

func (hf *HeaderFrame) FrameType() byte {
        return 2
}

func (hf *HeaderFrame) Read(reader io.Reader) error {
        class, err := ReadShort(reader)
        if err != nil {
                return err
        }
        hf.Class = class

        bodysize, err := ReadLongLong(reader)
        if err != nil {
                return err
        }
        hf.BodySize = bodysize
        return nil
}