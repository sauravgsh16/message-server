package proto

import (
	"io"
)

type Frame interface {
	Write(io.Writer) error
	Channel() uint16
}

type MessageFrame interface {
	MethodIdentifier() (uint16, uint16)
	Read(r io.Reader) (err error)
	Write(w io.Writer) (err error)
	FrameType() byte
	Wait() bool
	MethodName() string
}

type MessageContentFrame interface {
	MessageFrame
	GetBody() []byte
	SetBody([]byte)
}

type MessageResourceHolder interface {
	AcquireResources(qm *QueueMessage) bool
	ReleaseResources(qm *QueueMessage)
}
