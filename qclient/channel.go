package qclient

import (
	"bytes"
	"sync"

	"github.com/sauravgsh16/secoc-third/proto"
)

const (
	CH_INIT = iota
	CH_OPEN
	CH_CLOSING
	CH_CLOSED
)

type Channel struct {
	id             uint16
	incoming       chan *proto.WireFrame
	outgoing       chan *proto.WireFrame
	conn           *Connection
	consumers      *Consumer
	sendMux        sync.Mutex
	state          uint8
	currentMessage proto.MethodFrame
	header         *proto.HeaderFrame
	body           []byte
	errors         chan *proto.Error
}

func newChannel(c *Connection, id uint16) *Channel {
	return &Channel{
		id:        id,
		conn:      c,
		incoming:  make(chan *proto.WireFrame),
		outgoing:  c.outgoing,
		consumers: createConsumers(),
	}
}

func (ch *Channel) sendMethod(mf proto.MethodFrame) error {
	ch.sendMux.Lock()
	defer ch.sendMux.Unlock()

	if err := ch.conn.send(&proto.ChannelFrame{
		ChannelID: ch.id,
		Method:    mf,
	}); err != nil {
		return err
	}
	return nil
}

func (ch *Channel) sendContent(msg *proto.Message, mf proto.MethodFrame) error {
	ch.sendMux.Lock()
	defer ch.sendMux.Unlock()

	buf := bytes.NewBuffer(make([]byte, 0))

	// Write Headers
	proto.WriteShort(buf, msg.Header.Class)
	proto.WriteLongLong(buf, msg.Header.BodySize)

	// Send method
	if err := ch.sendMethod(mf); err != nil {
		return err
	}

	// Send Header
	ch.outgoing <- &proto.WireFrame{
		FrameType: uint8(proto.FrameHeader),
		Channel:   ch.id,
		Payload:   buf.Bytes(),
	}

	// Send Body
	for _, body := range msg.Payload {
		body.Channel = ch.id
		ch.outgoing <- body
	}
	return nil
}

func (ch *Channel) sendError(err *proto.Error) {
	if err.Soft {
		ch.sendMethod(&proto.ChannelClose{
			ReplyCode: err.Code,
			ReplyText: err.Msg,
			ClassId:   err.Class,
			MethodId:  err.Method,
		})
	} else {
		ch.conn.closeWithErr(err)
	}
}

func (ch *Channel) open() error {
	ch.startReceiver()
	if err := ch.sendMethod(&proto.ChannelOpen{}); err != nil {
		return err
	}
	return nil
}

func (ch *Channel) startReceiver() {
	if ch.state == 0 {
		ch.state = CH_OPEN
	}
	go func() {
		for {
			if ch.state == CH_CLOSED {
				break
			}
			var err *proto.Error
			frame := <-ch.incoming

			switch frame.FrameType {
			case uint8(proto.FrameMethod):
				err = ch.dispatchMethod(frame)
			case uint8(proto.FrameHeader):
				if ch.state != CH_CLOSING {
					err = ch.handleHeader(frame)
				}
			case uint8(proto.FrameBody):
				if ch.state != CH_CLOSING {
					err = ch.handleBody(frame)
				}
			default:
				err = proto.NewHardError(500, "Unknown frame type", 0, 0)
			}
			if err != nil {
				ch.sendError(err)
			}
		}

	}()
}

func (ch *Channel) dispatchMethod(wf *proto.WireFrame) *proto.Error {
	reader := bytes.NewReader(wf.Payload)

	mf, err := proto.ReadMethod(reader)
	if err != nil {
		return proto.NewHardError(500, err.Error(), 0, 0)
	}

	clsID, mtdID := mf.MethodIdentifier()

	switch clsID {
	case 10:
		return ch.connectionRoute(mf)
	case 20:
		return ch.channelRoute(mf)
	case 30:
		return ch.exchangeRoute(mf)
	case 40:
		return ch.queueRoute(mf)
	case 50:
		return ch.basicRoute(mf)
	default:
		return proto.NewHardError(500, "Unexpected method frame", clsID, mtdID)
	}
}

func (ch *Channel) handleHeader() *proto.Error {

}
