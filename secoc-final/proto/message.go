package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type WireFrame struct {
	FrameType uint8
	Channel   uint16
	Payload   []byte
}

// Used
type MethodFrame struct {
	ChannelID uint16
	ClassID   uint16
	MethodID  uint16
	Method    MessageFrame
}

func (mf *MethodFrame) Channel() uint16 { return mf.ChannelID }

func (mf *MethodFrame) Write(w io.Writer) error {
	var payload bytes.Buffer

	if mf.Method == nil {
		return errors.New("Missing Method - incorrectly frame")
	}

	clsID, mtdId := mf.Method.MethodIdentifier()

	if err := binary.Write(&payload, binary.BigEndian, clsID); err != nil {
		return err
	}

	if err := binary.Write(&payload, binary.BigEndian, mtdId); err != nil {
		return err
	}

	if err := mf.Method.Write(&payload); err != nil {
		return err
	}

	return writeFrame(w, FrameMethod, mf.ChannelID, payload.Bytes())
}

type HeaderFrame struct {
	Class     uint16
	ChannelID uint16
	BodySize  uint64
}

func (hf *HeaderFrame) Channel() uint16 { return hf.ChannelID }

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

type BodyFrame struct {
	ChannelID uint16
	Body      []byte
}

func (bf *BodyFrame) Channel() uint16 { return bf.ChannelID }

func (bf *BodyFrame) Write(w io.Writer) error {
	return writeFrame(w, FrameBody, bf.ChannelID, bf.Body)
}

type ChannelFrame struct {
	ChannelID uint16
	Method    MessageFrame
}

type Message struct {
	ID         int64
	Header     *HeaderFrame
	Payload    []*WireFrame
	Exchange   string
	RoutingKey string
	Method     MessageContentFrame
}

type QueueMessage struct {
	ID            int64
	DeliveryCount int32
	MsgSize       uint32
}

type TxMessage struct {
	Msg       *Message
	QueueName string
}

type IndexMessage struct {
	ID            int64
	Refs          int32
	DeliveryCount int32
	Persisted     bool
}

/*
func NewMessage(mcf MethodContentFrame) *Message {
	msg := &Message{
		ID:      NextCnt(),
		Payload: make([]*WireFrame, 0, 1),
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
*/

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
