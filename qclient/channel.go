package qclient

import (
	"sync"

	"github.com/sauravgsh16/secoc-third/proto"
)

type Channel struct {
	id             uint16
	incoming       chan *proto.MethodFrame
	outgoing       chan *proto.MethodFrame
	conn           *Connection
	consumers      *Consumer
	sendMux        sync.Mutex
	state          uint8
	currentMessage proto.MethodFrame
	header         *proto.HeaderFrame
	body           []byte
	recv           func(*Channel, *proto.WireFrame)
}

func newChannel(c *Connection, id uint16) *Channel {
	return &Channel{
		id:        id,
		conn:      c,
		incoming:  make(chan *proto.MethodFrame),
		outgoing:  make(chan *proto.MethodFrame),
		consumers: createConsumers(),
		recv:      (*Channel).recvFunc,
	}
}

func (ch *Channel) call(req proto.MethodFrame, res ...proto.MethodFrame) error {

	return nil
}

func (ch *Channel) open() error {
	return ch.call(&proto.ChannelOpen{}, &proto.ChannelOpenOk{})
}

func (ch *Channel) recvFunc(wf *proto.WireFrame) {}
