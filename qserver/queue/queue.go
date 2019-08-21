package queue

import (
	"errors"
	"sync"

	"github.com/sauravgsh16/secoc-third/qserver/consumer"
	"github.com/sauravgsh16/secoc-third/qserver/proto"
)

type Queue struct {
	Name        string
	list        *List
	Closed      bool
	qmux        sync.Mutex
	consumers   []*consumer.Consumer
	consumerMux sync.RWMutex
	ConnId      int64
	deleteChan  chan *Queue
}

func NewQueue(name string, connId int64, deleteChan chan *Queue) *Queue {
	return &Queue{
		Name:       name,
		list:       newlist(),
		consumers:  make([]*consumer.Consumer, 0),
		deleteChan: deleteChan,
	}
}

func (q *Queue) Len() uint32 {
	l := q.list.Len()
	if l < 0 {
		panic("Queue overflow")
	}
	return uint32(l)
}

func (q *Queue) ConsumerCount() uint32 {
	return uint32(len(q.consumers))
}

// ************************
//       IMPLEMENT
//

func (q *Queue) Close() {
	q.qmux.Lock()
	defer q.qmux.Unlock()
	q.Closed = true
}

func (q *Queue) Add(qm *proto.QueueMessage) bool {
	return true
}

func (q *Queue) Delete(ifUnused bool, ifEmpty bool) (uint32, error) {
	if !q.Closed {
		panic("Tryin to delete unclosed Queue")
	}
	q.qmux.Lock()
	defer q.qmux.Unlock()

	// Check if queue is being used
	used := !ifUnused || len(q.consumers) == 0
	emptied := !ifEmpty || q.list.Len() == 0

	if !used {
		return 0, errors.New("consumers present - specified unused")
	}
	if !emptied {
		return 0, errors.New("messages in Queue - specified is empty")
	}

	// Send cancel to all consumers of queue
	q.sendCancelConsumers()
	// Purge queue data
	return q.purgeQueueData(), nil
}

func (q *Queue) sendCancelConsumers() {
	q.consumerMux.Lock()
	defer q.consumerMux.Unlock()

	for _, c := range q.consumers {
		// ************* TODO *******************
		// Commenting below line
		// c.SendCancel()
		// **************************************
		c.Stop()
	}
	q.consumers = make([]*consumer.Consumer, 0, 1)
}

func (q *Queue) purgeQueueData() uint32 {
	length := q.list.Len()
	q.list.removeRef()
	return uint32(length)
}

func (q *Queue) GetOne(mrh ...proto.MessageResourceHolder) (*proto.QueueMessage, *proto.Message) {
	return nil, nil
}

func (q *Queue) AddConsumer(c *consumer.Consumer) (uint16, error) {
	return 0, nil
}
