package server

import (
	"fmt"

	"github.com/sauravgsh16/secoc-third/qserver/binding"

	"github.com/sauravgsh16/secoc-third/qserver/proto"
	"github.com/sauravgsh16/secoc-third/qserver/queue"
)

func (ch *Channel) queueRoute(mf proto.MethodFrame) *proto.ProtoError {
	switch method := mf.(type) {
	case *proto.QueueDeclare:
		return ch.queueDeclare(method)
	case *proto.QueueBind:
		return ch.queueBind(method)
	case *proto.QueueUnbind:
		return ch.queueUnbind(method)
	case *proto.QueueDelete:
		return ch.queueDelete(method)
	}
	clsID, mtdID := mf.MethodIdentifier()
	return proto.NewHardError(540, "Not Implemented Queue method", clsID, mtdID)
}

func (ch *Channel) queueDeclare(m *proto.QueueDeclare) *proto.ProtoError {
	clsID, mtdID := m.MethodIdentifier()

	q, ok := ch.conn.server.queues[m.Queue]
	if ok {
		qsize := uint32(q.Len())
		csize := q.ConsumerCount()
		ch.SendMethod(&proto.QueueDeclareOk{
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
		ch.SendMethod(&proto.QueueDeclareOk{
			Queue:       q.Name,
			MessageCnt:  uint32(0),
			ConsumerCnt: uint32(0),
		})
	}
	return nil
}

func (ch *Channel) queueBind(m *proto.QueueBind) *proto.ProtoError {
	clsID, mtdID := m.MethodIdentifier()

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
		ch.SendMethod(&proto.QueueBindOk{})
	}
	return nil
}

func (ch *Channel) queueUnbind(m *proto.QueueUnbind) *proto.ProtoError {
	clsID, mtdID := m.MethodIdentifier()

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

	ch.SendMethod(&proto.QueueUnbindOk{})
	return nil
}

func (ch *Channel) queueDelete(m *proto.QueueDelete) *proto.ProtoError {
	clsID, mtdID := m.MethodIdentifier()

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
		ch.SendMethod(&proto.QueueDeleteOk{MessageCnt: msgPurged})
	}
	return nil
}
