package queue

import (
	"errors"
        "sync"
        
        sh "github.com/sauravgsh16/secoc-third/shared"
)

//ErrEmptyQueue empty Queue
var ErrEmptyQueue = errors.New("attempt to get message from empty queue")

type Queue struct {
	ID   int   // Queue id for reference
	list *List
	mux  sync.Mutex
}

func NewQueue(id int) *Queue {
	l := newList()
	q := &Queue{
		ID  : id,
		list: l,
	}
	return q
}

// TODO: NEED TO IMPLEMENT MUTEX FOR WRITING TO AND READING FROM QUEUE
// EnQueue
func (q *Queue) EnQueue(d sh.QData) {
	q.mux.Lock()
	q.list.Append(d)
	q.mux.Unlock()
}

// DeQueue from queue
func (q *Queue) DeQueue() (sh.QData, error) {
	q.mux.Lock()
	d := q.list.Remove()
	q.mux.Unlock()
	if d == (sh.QData{}) {
		return sh.QData{}, ErrEmptyQueue
	}
	return d, nil
}

func (q *Queue) Len() int {
	return q.list.Len()
}
