package consumer

import (
	"sync"

	"github.com/sauravgsh16/secoc-third/qserver/proto"
	"github.com/sauravgsh16/secoc-third/qserver/store"
)

type Consumer struct {
	msgStore    *store.MsgStore
	ConsumerTag string
	chResource  ChannelResource
	incoming    chan bool
	cQueue      ConsumerQueue
	queueName   string
	mux         sync.Mutex
	stopped     bool
	stopMux     sync.Mutex
	noAck       bool
	defaultSize uint32
	activeSize  uint32
	sizeMux     sync.Mutex
}

type ConsumerQueue interface {
	GetOne(mrh ...proto.MessageResourceHolder) (*proto.QueueMessage, *proto.Message)
}

type ChannelResource interface {
	proto.MessageResourceHolder
	SendContent(mf proto.MethodFrame, msg *proto.Message)
	SendMethod(mf proto.MethodFrame)
	FlowActive() bool
	GetDeliveryTag() uint64
}

func NewConsumer(ms *store.MsgStore, cr ChannelResource, consumerTag string, cq ConsumerQueue, queueName string, noAck bool, defaultSize uint32) *Consumer {
	return &Consumer{
		msgStore:    ms,
		ConsumerTag: consumerTag,
		chResource:  cr,
		incoming:    make(chan bool),
		cQueue:      cq,
		queueName:   queueName,
		noAck:       noAck,
		defaultSize: defaultSize,
	}
}

func (c *Consumer) Start() {
	go c.consume()
}

func (c *Consumer) Stop() {
	c.stopMux.Lock()
	defer c.stopMux.Unlock()

	if !c.stopped {
		close(c.incoming)
		c.stopped = true
	}
}

func (c *Consumer) Ping() {
	c.stopMux.Lock()
	defer c.stopMux.Unlock()

	if !c.stopped {
		select {
		case c.incoming <- true:
		default:
		}
	}
}

func (c *Consumer) AcquireResources(qm *proto.QueueMessage) bool {

	c.sizeMux.Lock()
	defer c.sizeMux.Unlock()

	// If channel is already in use -
	// Client should use separate channels to publish and consume
	if !c.chResource.FlowActive() {
		return false
	}

	if c.noAck {
		c.activeSize += qm.MsgSize
		return true
	}

	if c.activeSize < c.defaultSize {
		c.activeSize += qm.MsgSize
		return true
	}
	return false
}

func (c *Consumer) ReleaseResources(qm *proto.QueueMessage) {

	c.sizeMux.Lock()
	defer c.sizeMux.Unlock()

	c.activeSize -= qm.MsgSize
}

func (c *Consumer) SendCancel() {
	c.chResource.SendMethod(&proto.BasicCancel{
		ConsumerTag: c.ConsumerTag,
		NoWait:      true,
	})
}

func (c *Consumer) ConsumeImmediate(msg *proto.Message, qm *proto.QueueMessage) bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	var tag uint64 = 0

	// TODO:
	/*
		if !c.noAck {
			tag := c.chResource.ADDUNACKMESSAGE()
		}
	*/
	c.chResource.SendContent(&proto.BasicDeliver{
		ConsumerTag: c.ConsumerTag,
		DeliveryTag: tag,
		Exchange:    msg.Exchange,
		RoutingKey:  msg.RoutingKey,
	}, msg)
	return true
}

func (c *Consumer) ResourceHolders() []proto.MessageResourceHolder {
	return []proto.MessageResourceHolder{c, c.chResource}
}

func (c *Consumer) consume() {
	for range c.incoming {
		c.consumeOne()
	}
}

func (c *Consumer) consumeOne() {
	c.mux.Lock()
	defer c.mux.Unlock()

	deliveryTag := c.chResource.GetDeliveryTag()

	_, msg := c.cQueue.GetOne(c.chResource, c)
	/*
		TODO
		if !c.noAck {
			We need to add this to a list of messages which have not
			been acknowledged yet.
		} else {
			We remove the reference of this message from the msg store
			as we will not see this message anymore.
		}
	*/
	c.chResource.SendContent(&proto.BasicDeliver{
		ConsumerTag: c.ConsumerTag,
		DeliveryTag: deliveryTag,
		Exchange:    msg.Exchange,
		RoutingKey:  msg.RoutingKey,
	}, msg)

	c.Ping()
}
