package server

import (
	"fmt"

	"github.com/sauravgsh16/secoc-third/proto"
	"github.com/sauravgsh16/secoc-third/qserver/binding"
	"github.com/sauravgsh16/secoc-third/qserver/queue"
)

func (ch *Channel) queueRoute(msgf proto.MessageFrame) *proto.Error {
	switch m := msgf.(type) {

	case *proto.QueueDeclare:
		return ch.qDeclare(m)

	case *proto.QueueBind:
		return ch.qBind(m)

	case *proto.QueueUnbind:
		return ch.qUnbind(m)

	case *proto.QueueDelete:
		return ch.qDelete(m)

	default:
		clsID, mtdID := msgf.Identifier()
		return proto.NewHardError(540, "Not Implemented Queue method", clsID, mtdID)
	}
}

func (ch *Channel) qDeclare(m *proto.QueueDeclare) *proto.Error {
	clsID, mtdID := m.Identifier()

	// Check if Queue already exists
	q, found := ch.conn.server.queues[m.Queue]
	if found {
		qsize := uint32(q.Len())
		csize := q.ConsumerCount()
		ch.Send(&proto.QueueDeclareOk{
			Queue:       m.Queue,
			MessageCnt:  qsize,
			ConsumerCnt: csize,
		})
		ch.usedQueueName = m.Queue
		return nil
	}

	// Create new Queue
	q = queue.NewQueue(m.Queue, ch.conn.id, ch.server.queueDeleter, ch.server.msgStore)
	// Add Queue
	err := ch.server.addQueue(q)
	if err != nil {
		return proto.NewSoftError(500, "failed to create a new Queue", clsID, mtdID)
	}
	ch.usedQueueName = q.Name
	if !m.NoWait {
		ch.Send(&proto.QueueDeclareOk{
			Queue:       q.Name,
			MessageCnt:  uint32(0),
			ConsumerCnt: uint32(0),
		})
	}
	return nil
}

func (ch *Channel) qBind(m *proto.QueueBind) *proto.Error {
	clsID, mtdID := m.Identifier()

	if len(m.Queue) == 0 {
		if len(ch.usedQueueName) == 0 {
			return proto.NewSoftError(404, "Queue not found", clsID, mtdID)
		}
		m.Queue = ch.usedQueueName
	}

	// Check queue
	q, found := ch.server.getQueue(m.Queue)
	if !found || q.Closed {
		return proto.NewSoftError(404, fmt.Sprintf("Queue: %s - not found", m.Queue), clsID, mtdID)
	}

	// Exchange queue
	ex, found := ch.server.getExchange(m.Exchange)
	if !found {
		return proto.NewSoftError(404, "Exchange not found", clsID, mtdID)
	}

	if q.ConnId != -1 && q.ConnId != ch.conn.id {
		return proto.NewSoftError(410, "Queue is binded to a different connection", clsID, mtdID)
	}

	// Create binding
	b, err := binding.NewBinding(m.Queue, m.Exchange, m.RoutingKey)
	if err != nil {
		return proto.NewSoftError(500, err.Error(), clsID, mtdID)
	}

	// Add the binding to the exchange
	err = ex.AddBinding(b, ch.conn.id)
	if err != nil {
		return proto.NewSoftError(500, err.Error(), clsID, mtdID)
	}

	if !m.NoWait {
		ch.Send(&proto.QueueBindOk{})
	}
	return nil
}

func (ch *Channel) qUnbind(m *proto.QueueUnbind) *proto.Error {
	clsID, mtdID := m.Identifier()

	if len(m.Queue) == 0 {
		if len(ch.usedQueueName) == 0 {
			return proto.NewSoftError(404, "Queue not found", clsID, mtdID)
		}
		m.Queue = ch.usedQueueName
	}

	// Check queue
	q, found := ch.server.queues[m.Queue]
	if !found || q.Closed {
		return proto.NewSoftError(404, fmt.Sprintf("Queue: %s - not found", m.Queue), clsID, mtdID)
	}

	// Exchange queue
	ex, found := ch.server.exchanges[m.Exchange]
	if !found {
		return proto.NewSoftError(404, "Exchange not found", clsID, mtdID)
	}

	if q.ConnId != -1 && q.ConnId != ch.conn.id {
		return proto.NewSoftError(410, "Queue is binded to a different connection", clsID, mtdID)
	}

	binding, err := binding.NewBinding(m.Queue, m.Exchange, m.RoutingKey)
	if err != nil {
		return proto.NewSoftError(500, err.Error(), clsID, mtdID)
	}

	err = ex.RemoveBinding(binding)
	if err != nil {
		return proto.NewSoftError(500, err.Error(), clsID, mtdID)
	}

	ch.Send(&proto.QueueUnbindOk{})
	return nil
}

func (ch *Channel) qDelete(m *proto.QueueDelete) *proto.Error {
	clsID, mtdID := m.Identifier()

	if len(m.Queue) == 0 {
		if len(ch.usedQueueName) == 0 {
			return proto.NewSoftError(404, "Queue not found", clsID, mtdID)
		}
		m.Queue = ch.usedQueueName
	}

	msgPurged, errCode, err := ch.server.deleteQueue(m, ch.conn.id)
	if err != nil {
		return proto.NewSoftError(errCode, err.Error(), clsID, mtdID)
	}

	if !m.NoWait {
		ch.Send(&proto.QueueDeleteOk{MessageCnt: msgPurged})
	}
	return nil
}
