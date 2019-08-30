package qclient

import (
	"reflect"
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

func (ch *Channel) send() error {

}

func (ch *Channel) call(req proto.MethodFrame, resps ...proto.MethodFrame) error {
	if err := ch.send(); err != nil {
		return err
	}

	if req.Wait() {
		select {
		case e, ok := <-ch.errors:
			if ok {
				return e
			}
			return ErrClosed
		case msg := <-ch.incoming:
			for _, res := range resps {
				if reflect.TypeOf(msg) == reflect.TypeOf(res) {
					vres := reflect.ValueOf(res).Elem()
					vmsg := reflect.ValueOf(msg).Elem()
					vres.Set(vmsg)
					return nil
				}
			}
			return ErrInvalidCommand
		}
	}
	return nil
}

func (ch *Channel) open() error {
	ch.startReceiver()
	return ch.openChannel()
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
				err = ch.routeMethod(frame)
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

func (ch *Channel) openChannel() error {
	return ch.call(&proto.ChannelOpen{}, &proto.ChannelOpenOk{})
}
