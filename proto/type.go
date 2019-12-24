package proto

import (
	"io"
)

// Frame interface
// Defines a Write method and a channel method
type Frame interface {
	Write(io.Writer) error
	Channel() uint16
}

// MessageFrame interface defines the interface for a message frame
// These are instructions which form the backbone for the message transfer strategies
type MessageFrame interface {
	Identifier() (uint16, uint16)
	Read(r io.Reader) (err error)
	Write(w io.Writer) (err error)
	FrameType() byte
	Wait() bool
	MethodName() string
}

// MessageContentFrame inteface defines an inteface for message
// frames which contains actual content which gets transferred
type MessageContentFrame interface {
	MessageFrame
	GetBody() []byte
	SetBody([]byte)
}

// MessageResourceHolder interface defines the interface for
// responsible for holding the resource
type MessageResourceHolder interface {
	AcquireResources(qm *QueueMessage) bool
	ReleaseResources(qm *QueueMessage)
}
