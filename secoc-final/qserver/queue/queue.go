package queue

import (
	"errors"
	"fmt"
	"sync"

	"github.com/sauravgsh16/secoc-third/secoc-final/proto"
	"github.com/sauravgsh16/secoc-third/secoc-final/qserver/consumer"
	"github.com/sauravgsh16/secoc-third/secoc-final/qserver/store"
)

type Queue struct {
	Name               string
	list               *List
	Closed             bool
	mux                sync.Mutex
	consumers          []*consumer.Consumer
	consumerMux        sync.RWMutex
	ConnId             int64
	deleteChan         chan *Queue
	readyChan          chan bool
	currentConsumerIdx int
	msgStore           *store.MsgStore
}

func NewQueue(name string, connId int64, deleteChan chan *Queue, msgStore *store.MsgStore) *Queue {
	return &Queue{
		Name:       name,
		list:       newlist(),
		consumers:  make([]*consumer.Consumer, 0, 1),
		deleteChan: deleteChan,
		readyChan:  make(chan bool, 1),
		msgStore:   msgStore,
	}
}

func (q *Queue) Start() {
	go func() {
		select {
		case q.readyChan <- true:
		default:
		}
		for _ = range q.readyChan {
			if q.Closed {
				fmt.Printf("Queue Closed: %s\n", q.Name)
				break
			}
			q.processSingleEntry()
		}
	}()
}

func (q *Queue) Len() uint32 {
	l := q.list.Len()
	if l < 0 {
		panic("Queue overflow")
	}
	return uint32(l)
}

func (q *Queue) Close() {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.Closed = true
}

func (q *Queue) Add(qm *proto.QueueMessage) bool {
	q.mux.Lock()
	defer q.mux.Unlock()

	if q.Closed {
		return false
	}
	q.list.Append(qm)
	select {
	case q.readyChan <- true:
	default:
	}
	return true
}

func (q *Queue) Delete(ifUnused bool, ifEmpty bool) (uint32, error) {
	if !q.Closed {
		panic("Tryin to delete unclosed Queue")
	}
	q.mux.Lock()
	defer q.mux.Unlock()

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
	q.cancelConsumers()
	// Purge queue data
	return q.purgeQueueData(), nil
}

func (q *Queue) ConsumeImmediate(qm *proto.QueueMessage) bool {
	q.consumerMux.Lock()
	defer q.consumerMux.Unlock()

	for _, consumer := range q.consumers {
		msg, acquired := q.msgStore.Get(qm, consumer.ResourceHolders())
		if acquired {
			return consumer.ConsumeImmediate(msg, qm)
		}
	}
	return false
}

func (q *Queue) processSingleEntry() {
	q.consumerMux.Lock()
	defer q.consumerMux.Unlock()

	length := len(q.consumers)
	if length == 0 {
		return
	}
	for count := 0; count < length; count++ {
		q.currentConsumerIdx = (q.currentConsumerIdx + 1) % length
		c := q.consumers[q.currentConsumerIdx]
		c.Ping()
	}
}

func (q *Queue) ConsumerCount() uint32 {
	return uint32(len(q.consumers))
}

func (q *Queue) AddConsumer(c *consumer.Consumer) (uint16, error) {
	if q.Closed {
		return 0, nil // Check: if the error should be nil here ?
	}

	q.consumerMux.Lock()
	q.consumers = append(q.consumers, c)
	q.consumerMux.Unlock()
	return 0, nil // Check: if it really needs to return the number of consumers ?
}

func (q *Queue) cancelConsumers() {
	q.consumerMux.Lock()
	defer q.consumerMux.Unlock()

	for _, c := range q.consumers {
		c.SendCancel()
		c.Stop()
	}
	q.consumers = make([]*consumer.Consumer, 0, 1)
}

func (q *Queue) removeConsumers(consumerTag string) {
	q.consumerMux.Lock()
	defer q.consumerMux.Unlock()

	// Remove consumers based on consumerTag
	for i, c := range q.consumers {
		if c.ConsumerTag == consumerTag {
			q.consumers = append(q.consumers[:i], q.consumers[i+1:]...)
		}
	}
	if len(q.consumers) == 0 {
		q.currentConsumerIdx = 0
	} else {
		q.currentConsumerIdx = q.currentConsumerIdx % len(q.consumers)
	}
}

func (q *Queue) purgeQueueData() uint32 {
	length := q.list.Len()
	q.list.removeRef()
	return uint32(length)
}

func (q *Queue) GetOne(mrh ...proto.MessageResourceHolder) (*proto.QueueMessage, *proto.Message) {
	q.mux.Lock()
	defer q.mux.Unlock()

	if q.Closed || q.Len() == 0 {
		return nil, nil
	}

	qm := q.list.Front().(*proto.QueueMessage)
	if qm == nil {
		return nil, nil
	}
	msg, acquired := q.msgStore.Get(qm, mrh)
	if !acquired {
		return nil, nil
	}
	q.list.Remove()
	return qm, msg
}
