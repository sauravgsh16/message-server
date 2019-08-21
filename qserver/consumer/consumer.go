package consumer

import (
	"sync"

	"github.com/sauravgsh16/secoc-third/qserver/proto"
	"github.com/sauravgsh16/secoc-third/qserver/store"
)

type Consumer struct {
	msgStore    *store.MsgStore
	ConsumerTag string
	cResource   ConsumerResource
	incoming    chan bool
	cQueue      ConsumerQueue
	queueName   string
	mux         sync.Mutex
	stopped     bool
	stopMux     sync.Mutex
	noAck       bool
}

type ConsumerQueue interface {
	GetOne(mrh ...proto.MessageResourceHolder) (*proto.QueueMessage, *proto.Message)
}

type ConsumerResource interface {
	SendContent(mf proto.MethodFrame, msg *proto.Message)
	SendMethod(mf proto.MethodFrame)
	FlowActive() bool
	GetDeliveryTag() uint64
}

func NewConsumer(ms *store.MsgStore, cr ConsumerResource, consumerTag string, cq ConsumerQueue, queueName string, noAck bool) *Consumer {
	return &Consumer{
		msgStore:    ms,
		ConsumerTag: consumerTag,
		cResource:   cr,
		incoming:    make(chan bool),
		cQueue:      cq,
		queueName:   queueName,
		noAck:       noAck,
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

func (c *Consumer) consume() {
	for range c.incoming {
		c.consumeOne()
	}
}

func (c *Consumer) consumeOne() {
	c.mux.Lock()
	defer c.mux.Unlock()
	deliveryTag := c.cResource.GetDeliveryTag()
	_, msg := c.cQueue.GetOne(c.cResource, c)
	c.cResource.SendContent(&proto.BasicDeliver{
		ConsumerTag: c.ConsumerTag,
		DeliveryTag: deliveryTag,
		Exchange:    msg.Exchange,
		RoutingKey:  msg.RoutingKey,
	}, msg)

	c.Ping()
}
