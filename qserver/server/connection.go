package server

import (
        "net"
        "sync"
        "sync/atomic"
        "math/rand"
        "time"

        "github.com/sauravgsh16/secoc-third/allocate"
        "github.com/sauravgsh16/secoc-third/qserver/proto"
)

// Below struct to have their own files

// ***** END *********

var counter int64

func init() {
        rand.Seed(time.Now().UnixNano())
        counter = time.Now().UnixNano()
}

func nextId() int64 {
        return atomic.AddInt64(&counter, 1)
}

type Connection struct {
        id        int64
        channels  map[uint64]*Channel
        outgoing  chan *proto.WireFrame
        server    *Server
        network   net.Conn
        mux       sync.Mutex
        allocator allocate.Allocator
}

func NewConnection(s *Server, n net.Conn) *Connection {
        return &Connection{
                id:       nextId(),
                channels: make(map[uint64]*Channel),
                outgoing: make(chan *proto.WireFrame, 100),
                server:   s,
                network:  n,       
        }
}

func (conn *Connection) openConnection() {
        // Create channel 0 and start the connection handshake
        conn.channels[0] = NewChannel(0, conn)
}