package server

import (
        "fmt"
        "sync"
        "net"

        "github.com/sauravgsh16/secoc-third/qserver/exchange"
        "github.com/sauravgsh16/secoc-third/qserver/queue"
        "github.com/sauravgsh16/secoc-third/qserver/proto"
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

func (s *Server) OpenConnection(conn net.Conn) {
        c := NewConnection(s, conn)
        s.conns[c.id] = c
        c.openConnection()
}