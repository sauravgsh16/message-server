package queue

import (
	"sync"

	"github.com/sauravgsh16/secoc-third/qserver/consumer"
	sh "github.com/sauravgsh16/secoc-third/shared"
)

// qData : need to check if more implementation details required
// Check Table : - types.go amqp
type qData struct {
	data sh.Message
}

type Queue struct {
        Name         string
	list         *List
	mux          sync.Mutex
        In           chan sh.Message
	Out          chan sh.Message
	consumers    []*consumer.Consumer
	ConnId       int64
	deleteChan   chan *Queue
	qLock        sync.Mutex
	consumerLock sync.RWMutex
}

func NewQueue(name string, connId int64, deleteChan chan *Queue) *Queue {
	q := &Queue{
                Name: name,
                In:   make(chan sh.Message),
		Out:  make(chan sh.Message),
		consumers:  make([]*consumer.Consumer, 0),
		deleteChan: deleteChan,
	}
	go q.datapump()
	return q
}

func (q *Queue) ConsumerCount() uint32 {
	return uint32(len(q.consumers))
}

func (q *Queue) datapump() {
channel:
	for {
                // new queue
                if q.list == nil {
                        q.list = newList()
                }
		select {
		case msg, ok := <- q.In:
			if !ok {
				break channel // Input channel closed, we break out of loop
			}
			q.enQueue(msg)
		// When reading from Out channel, need to keep in mind
		// when no data is present, we are returning an empty Message struct
		case q.Out <- q.deQueue():
		}
	}
	// We drain the output channel
	for q.Len() > 0 {
		q.Out <- q.deQueue()
	}
	close(q.Out)
}

// TODO: NEED TO IMPLEMENT MUTEX FOR WRITING TO AND READING FROM QUEUE
// EnQueue
func (q *Queue) enQueue(msg sh.Message) {
	qd := qData{data: msg}
	q.mux.Lock()
	q.list.Append(qd)
	q.mux.Unlock()
}

// DeQueue from queue
func (q *Queue) deQueue() sh.Message {
	q.mux.Lock()
	qd := q.list.Remove()
	q.mux.Unlock()
	return qd.data
}

func (q *Queue) Len() int {
	return q.list.Len()
}