package server

import (
	"github.com/sauravgsh16/message-server/allocate"
	"github.com/sauravgsh16/message-server/proto"
)

func (ch *Channel) basicRoute(msgf proto.MessageFrame) *proto.Error {
	switch m := msgf.(type) {

	case *proto.BasicConsume:
		return ch.basicConsume(m)

	case *proto.BasicCancel:
		return ch.basicCancel(m)

	case *proto.BasicPublish:
		return ch.basicPublish(m)

	case *proto.BasicAck:
		// TODO: To implement
		return ch.basicAck(m)

	case *proto.BasicNack:
		// TODO: To implement
		return ch.basicNack(m)

	default:
		clsID, mtdID := msgf.Identifier()
		return proto.NewHardError(540, "unable to route method frame", clsID, mtdID)
	}
}

func (ch *Channel) basicConsume(m *proto.BasicConsume) *proto.Error {
	clsID, mtdID := m.Identifier()

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
		ch.Send(&proto.BasicConsumeOk{ConsumerTag: m.ConsumerTag})
	}
	return nil
}

func (ch *Channel) basicCancel(m *proto.BasicCancel) *proto.Error {
	if err := ch.removeConsumer(m.ConsumerTag); err != nil {
		clsID, mtdID := m.Identifier()
		return proto.NewSoftError(404, err.Error(), clsID, mtdID)
	}
	if !m.NoWait {
		ch.Send(&proto.BasicCancelOk{ConsumerTag: m.ConsumerTag})
	}
	return nil
}

func (ch *Channel) basicPublish(m *proto.BasicPublish) *proto.Error {
	_, found := ch.server.exchanges[m.Exchange]
	if !found {
		clsID, mtdID := m.Identifier()
		return proto.NewSoftError(404, "Exchange not found", clsID, mtdID)
	}

	ch.startPublish(m)
	return nil
}

func (ch *Channel) basicAck(m *proto.BasicAck) *proto.Error {
	return nil
}

func (ch *Channel) basicNack(m *proto.BasicNack) *proto.Error {
	return nil
}
