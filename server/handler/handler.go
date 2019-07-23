package handler

import (
	"fmt"
	"net/rpc"
	"sync"

	"github.com/sauravgsh16/secoc-third/server/queue"
	sh "github.com/sauravgsh16/secoc-third/shared"
)

type qref map[string]*queue.Queue

func (qr qref) add(n string, q *queue.Queue) {
	qr[n] = q
}

func (qr qref) get(n string) (*queue.Queue, error) {
	q, ok := qr[n]
	if !ok {
		return &queue.Queue{}, fmt.Errorf("queue with id: %s - not found", n)
	}
	return q, nil
}

/*
Store - stores all the information when a new instance of a server is create
*/
type Store struct {
	queues map[string]qref
	wmux   sync.Mutex
	rwmux  sync.RWMutex
}

func (ss *Store) initq(qtype string) qref {
	child, ok := ss.queues[qtype]
	if !ok {
		child = make(map[string]*queue.Queue)
		ss.queues[qtype] = child
	}
	return child
}

// registers new queue
func (ss *Store) newQ(qtype, qName string) {
	q := queue.NewQueue()
	qr := ss.initq(qtype)
	ss.wmux.Lock()
	qr.add(qName, q)
	ss.wmux.Unlock()
}

// getQ from store
func (ss *Store) getQ(qtype, qName string) (*queue.Queue, error) {
	qr := ss.queues[qtype]
	q, err := qr.get(qName)
	if err != nil {
		return &queue.Queue{}, fmt.Errorf("queue with name: %s - not found", qName)
	}
	return q, nil
}

/*
HandleQ - responsible for exposing all rpc methods for clients to call
*/
type HandleQ struct {
	ss *Store
}

func (h *HandleQ) validateQtype(qtype string) error {
	switch qtype {
	case "broadcast", "send":
	default:
		return fmt.Errorf("invalid Queue type: %s", qtype)
	}
	return nil
}

// CreateQueue - creates new queue into store
func (h *HandleQ) CreateQueue(args *sh.QCreate, res *int) error {
	err := h.validateQtype(args.Qtype)
	if err != nil {
		return err
	}
	h.ss.newQ(args.Qtype, args.QName)
	return nil
}

// Publish data to the queue
func (h *HandleQ) Publish(args *sh.QPublish, res *int) error {
	err := h.validateQtype(args.Qtype)
	if err != nil {
		return err
	}
	q, err := h.ss.getQ(args.Qtype, args.QName)
	if err != nil {
		return err
	}
	q.EnQueue(args.Qd)
	*res = 1
	return nil
}

// FetchData - from queue
func (h *HandleQ) FetchData(args *sh.QFetch, res *sh.QData) error {
	if err := h.validateQtype(args.Qtype); err != nil {
		return err
	}
	q, err := h.ss.getQ(args.Qtype, args.QName)
	if err != nil {
		return err
	}
	qd, err := q.DeQueue()
	if err != nil {
		return err
	}
	*res = qd
	return nil
}

// Register Handlers
func Register() error {
	ss := &Store{
		queues: make(map[string]qref),
	}
	h := HandleQ{ss}
	if err := rpc.Register(&h); err != nil {
		return err
	}
	return nil
}