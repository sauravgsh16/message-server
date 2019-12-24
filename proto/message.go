package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// ProtocolHeader struct represents the initial exchange between
// client and server for handshake
type ProtocolHeader struct{}

// Channel - not implemented
// Present to satisfy Frame interface
func (*ProtocolHeader) Channel() uint16 {
	panic("Should never be called")
}

func (*ProtocolHeader) Write(w io.Writer) (err error) {
	_, err = w.Write([]byte{'S', 'E', 'C', 'O', 'C'})
	return err
}

// MethodFrame struct represents the method call
type MethodFrame struct {
	ChannelID uint16
	ClassID   uint16
	MethodID  uint16
	Method    MessageFrame
}

// Channel returns the channel id
func (mf *MethodFrame) Channel() uint16 { return mf.ChannelID }

func (mf *MethodFrame) Write(w io.Writer) error {
	var payload bytes.Buffer

	if mf.Method == nil {
		return errors.New("Missing Method - incorrectly frame")
	}

	clsID, mtdID := mf.Method.Identifier()

	if err := binary.Write(&payload, binary.BigEndian, clsID); err != nil {
		return err
	}

	if err := binary.Write(&payload, binary.BigEndian, mtdID); err != nil {
		return err
	}

	if err := mf.Method.Write(&payload); err != nil {
		return err
	}

	return writeFrame(w, FrameMethod, mf.ChannelID, payload.Bytes())
}

// HeaderFrame struct represent the header frame of the message
type HeaderFrame struct {
	Class     uint16
	ChannelID uint16
	BodySize  uint64
}

// Channel returns the channel id
func (hf *HeaderFrame) Channel() uint16 { return hf.ChannelID }

// FrameType returns the frame type
func (hf *HeaderFrame) FrameType() byte { return 2 }

func (hf *HeaderFrame) Write(w io.Writer) error {
	var payload bytes.Buffer

	if err := binary.Write(&payload, binary.BigEndian, hf.Class); err != nil {
		return err
	}

	if err := binary.Write(&payload, binary.BigEndian, hf.BodySize); err != nil {
		return err
	}

	return writeFrame(w, FrameHeader, hf.ChannelID, payload.Bytes())
}

// BodyFrame struct contains the actual body of the message
type BodyFrame struct {
	ChannelID uint16
	Body      []byte
}

// Channel returns the channel id
func (bf *BodyFrame) Channel() uint16 { return bf.ChannelID }

func (bf *BodyFrame) Write(w io.Writer) error {
	return writeFrame(w, FrameBody, bf.ChannelID, bf.Body)
}

// Message struct contains the information required to send msg
// Contains MessageContentFrame
type Message struct {
	ID         int64
	Header     *HeaderFrame
	Payload    []byte
	Exchange   string
	RoutingKey string
	Method     MessageContentFrame
}

// QueueMessage struct
type QueueMessage struct {
	ID            int64
	DeliveryCount int32
	MsgSize       uint32
}

// TxMessage struct
type TxMessage struct {
	Msg       *Message
	QueueName string
}

// IndexMessage struct
type IndexMessage struct {
	ID            int64
	Refs          int32
	DeliveryCount int32
	Persisted     bool
}

// NewMessage returns a new message.
// Take MessageContentFrame as input
func NewMessage(mcf MessageContentFrame) *Message {
	msg := &Message{
		ID:      NextCnt(),
		Payload: make([]byte, 0),
	}
	switch m := mcf.(type) {
	case *BasicPublish:
		msg.Method = m
		msg.Exchange = m.Exchange
		msg.RoutingKey = m.RoutingKey
	case *BasicDeliver:
		msg.Method = m
		msg.Exchange = m.Exchange
		msg.RoutingKey = m.RoutingKey
	}
	return msg
}

// NewTxMessage returns a new TxMessage.
// Takes a pointer to Message and queue name as input
func NewTxMessage(msg *Message, qn string) *TxMessage {
	return &TxMessage{
		Msg:       msg,
		QueueName: qn,
	}
}

// NewIndexMessage returns a new IndexMessage.
// Takes index id, reference count and delivery count as input
func NewIndexMessage(id int64, refcount int32, deliveryCount int32) *IndexMessage {
	return &IndexMessage{
		ID:            id,
		Refs:          refcount,
		DeliveryCount: deliveryCount,
	}
}

// NewQueueMessage returns a new QueueMessage.
// Takes queue id, delivery count and message size as imput
func NewQueueMessage(id int64, deliveryCount int32, size uint32) *QueueMessage {
	return &QueueMessage{
		ID:            id,
		DeliveryCount: deliveryCount,
		MsgSize:       size,
	}
}
