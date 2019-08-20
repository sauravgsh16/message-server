package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/sauravgsh16/secoc-third/qserver/exchange"
	"github.com/sauravgsh16/secoc-third/qserver/proto"
	"github.com/sauravgsh16/secoc-third/qserver/queue"
)

type Server struct {
	exchanges       map[string]*exchange.Exchange
	queues          map[string]*queue.Queue
	conns           map[int64]*Connection
	mux             sync.Mutex
	exchangeDeleter chan *exchange.Exchange
	queueDeleter    chan *queue.Queue
}

func NewServer() *Server {
	var server = &Server{
		exchanges:       make(map[string]*exchange.Exchange),
		queues:          make(map[string]*queue.Queue),
		conns:           make(map[int64]*Connection),
		exchangeDeleter: make(chan *exchange.Exchange),
		queueDeleter:    make(chan *queue.Queue),
	}
	return server
}

func (s *Server) addExchange(ex *exchange.Exchange) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.exchanges[ex.Name] = ex
	return nil
}

func (s *Server) deleteExchange(m *proto.ExchangeDelete) (uint16, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	ex, ok := s.exchanges[m.Exchange]
	if !ok {
		return 404, fmt.Errorf("Exchange: %s - not found", m.Exchange)
	}
	// Close everything associated with the exchange
	ex.Close()
	delete(s.exchanges, m.Exchange)
	return 0, nil
}

func (s *Server) addQueue(q *queue.Queue) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.queues[q.Name] = q
	return nil
}

func (s *Server) deleteQueue(m *proto.QueueDelete, connID int64) (uint32, uint16, error) {

	s.mux.Lock()
	defer s.mux.Unlock()

	q, found := s.queues[m.Queue]
	if !found {
		return 0, 404, fmt.Errorf("Queue not found")
	}

	if q.ConnId != -1 && q.ConnId != connID {
		return 0, 405, fmt.Errorf("Queue is locked by another connection")
	}

	// Close queue - to stop any data enqueue and dequeue
	q.Close()
	// Remove queue from all the bindings
	s.removeQueueBindings(m.Queue)

	// Cleanup
	msgPurged, err := q.Delete(m.IfUnused, m.IfEmpty)
	if err != nil {
		return 0, 406, err
	}
	delete(s.queues, m.Queue)
	return msgPurged, 0, nil
}

func (s *Server) removeQueueBindings(qName string) {
	for _, ex := range s.exchanges {
		ex.RemoveQueueBindings(qName)
	}
}

func (s *Server) deleteConnection(connID int64) {
	delete(s.conns, connID)
}

func (s *Server) OpenConnection(conn net.Conn) {
	c := NewConnection(s, conn)
	s.conns[c.id] = c
	c.openConnection()
}
