package server

import (
        "sync"
        "net"

        "github.com/sauravgsh16/secoc-third/qserver/queue"
)

// Below struct to have their own files
type Exchange struct {
        Name string
}
// ***** END *********


type Server struct {
        exchanges map[string]*Exchange
        queues    map[string]*queue.Queue
        conns     map[int64]*Connection
        mux       sync.Mutex
}

func NewServer() *Server {
        var server = &Server{
                exchanges: make(map[string]*Exchange),
                queues:    make(map[string]*queue.Queue),
                conns:     make(map[int64]*Connection),
        }
        return server
}

func (s *Server) addExchange(ex *Exchange) error {
        s.mux.Lock()
        defer s.mux.Unlock()
        s.exchanges[ex.Name] = ex
        return nil
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