package server

import (
	"github.com/sauravgsh16/secoc-third/allocate"
	"github.com/sauravgsh16/secoc-third/qserver/proto"
)

func (ch *Channel) basicRoute(mf proto.MethodFrame) *proto.ProtoError {
	switch method := mf.(type) {
	case *proto.BasicConsume:
		return ch.basicConsume(method)
	case *proto.BasicCancel:
		return ch.basicCancel(method)
	case *proto.BasicPublish:
		return ch.basicPublish(method)
	case *proto.BasicAck:
		return ch.basicAck(method)
	case *proto.BasicNack:
		return ch.basicNack(method)
	default:
		clsID, mtdID := mf.MethodIdentifier()
		return proto.NewHardError(540, "unable to route method frame", clsID, mtdID)
	}
}

func (ch *Channel) basicConsume(m *proto.BasicConsume) *proto.ProtoError {
	clsID, mtdID := m.MethodIdentifier()

	// Check queue
	if len(m.Queue) == 0 {
		if len(ch.usedQueueName) == 0 {
			return proto.NewSoftError(404, "Queue not found", clsID, mtdID)
		}
		m.Queue = ch.usedQueueName
	}

	q, found := ch.server.queues[m.Queue]
	if !found {
		return proto.NewSoftError(404, "Queue not found", clsID, mtdID)
	}

	if len(m.ConsumerTag) == 0 {
		m.ConsumerTag = allocate.RandomID()
	}

	err := ch.addNewConsumer(q, m)
	if err != nil {
		return err
	}

	if !m.NoWait {
		ch.SendMethod(&proto.BasicConsumeOk{ConsumerTag: m.ConsumerTag})
	}
	return nil
}

func (ch *Channel) basicCancel(m *proto.BasicCancel) *proto.ProtoError {
	if err := ch.removeConsumer(m.ConsumerTag); err != nil {
		clsID, mtdID := m.MethodIdentifier()
		return proto.NewSoftError(404, err.Error(), clsID, mtdID)
	}
	if !m.NoWait {
		ch.SendMethod(&proto.BasicCancelOk{ConsumerTag: m.ConsumerTag})
	}
	return nil
}

func (ch *Channel) basicPublish(m *proto.BasicPublish) *proto.ProtoError {
	_, found := ch.server.exchanges[m.Exchange]
	if !found {
		clsID, mtdID := m.MethodIdentifier()
		return proto.NewSoftError(404, "Exchange not found", clsID, mtdID)
	}

	ch.startPublish(m)
	return nil
}

func (ch *Channel) basicAck(m *proto.BasicAck) *proto.ProtoError {
	return nil
}

func (ch *Channel) basicNack(m *proto.BasicNack) *proto.ProtoError {
	return nil
}
