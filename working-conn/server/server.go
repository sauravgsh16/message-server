package server

import "net"

type Server struct {
	conn map[int64]*Connection
}

func NewServer() *Server {
	s := &Server{
		conn: make(map[int64]*Connection),
	}
	return s
}

func (s *Server) OpenConnection(conn net.Conn) {
	c := NewConn(s, conn)
	s.conn[int64(1)] = c
	c.openConn()
}
