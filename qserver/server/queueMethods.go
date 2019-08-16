package server

import (
	"github.com/sauravgsh16/secoc-third/qserver/queue"
        "github.com/sauravgsh16/secoc-third/qserver/proto"
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
        clsId, mtdId := mf.MethodIdentifier()
        return proto.NewHardError(540, "Not Implemented Queue method", clsId, mtdId)
}

func (ch *Channel) queueDeclare(m *proto.QueueDeclare) *proto.ProtoError {
        clsId, mtdId := m.MethodIdentifier()

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
        q = queue.NewQueue(m.Queue, ch.conn.id, ch.server.queueDeleter)
        // Add Queue
        err := ch.server.addQueue(q)
        if err != nil {
                return proto.NewSoftError(500, "failed to create a new Queue", clsId, mtdId)
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
        return nil
}

func (ch *Channel) queueUnbind(m *proto.QueueUnbind) *proto.ProtoError {
        return nil
}

func (ch *Channel) queueDelete(m *proto.QueueDelete) *proto.ProtoError {
        return nil
}