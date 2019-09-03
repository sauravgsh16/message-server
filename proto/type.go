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
	Wait() bool
}

type MethodContentFrame interface {
	MethodFrame
	GetBody() []byte
	SetBody([]byte)
}

type MessageResourceHolder interface {
	AcquireResources(qm *QueueMessage) bool
	ReleaseResources(qm *QueueMessage)
}

func NewMessage(m *BasicPublish) *Message {
	return &Message{
		ID:         NextCnt(),
		Method:     m,
		Exchange:   m.Exchange,
		RoutingKey: m.RoutingKey,
		Payload:    make([]*WireFrame, 0, 1),
	}
}

func NewTxMessage(msg *Message, qn string) *TxMessage {
	return &TxMessage{
		Msg:       msg,
		QueueName: qn,
	}
}

func NewIndexMessage(id int64, refcount int32, deliveryCount int32) *IndexMessage {
	return &IndexMessage{
		ID:            id,
		Refs:          refcount,
		DeliveryCount: deliveryCount,
	}
}

func NewQueueMessage(id int64, deliveryCount int32, size uint32) *QueueMessage {
	return &QueueMessage{
		ID:            id,
		DeliveryCount: deliveryCount,
		MsgSize:       size,
	}
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
