package server

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sauravgsh16/secoc-third/allocate"
	"github.com/sauravgsh16/secoc-third/qserver/proto"
)

var counter int64

func init() {
	rand.Seed(time.Now().UnixNano())
	counter = time.Now().UnixNano()
}

func nextId() int64 {
	return atomic.AddInt64(&counter, 1)
}

type ConnectionStatus struct {
	start    bool
	startOk  bool
	open     bool
	openOk   bool
	closing  bool
	closed   bool
	closedOk bool
}

type Connection struct {
	id        int64
	channels  map[uint16]*Channel
	outgoing  chan *proto.WireFrame
	server    *Server
	network   net.Conn
	mux       sync.Mutex
	allocator allocate.Allocator
	status    ConnectionStatus
}

func NewConnection(s *Server, n net.Conn) *Connection {
	return &Connection{
		id:       nextId(),
		channels: make(map[uint16]*Channel),
		outgoing: make(chan *proto.WireFrame, 100),
		server:   s,
		network:  n,
		status:   ConnectionStatus{},
	}
}

func (conn *Connection) openConnection() {
	// Create channel 0 and start the connection handshake
	conn.channels[0] = NewChannel(0, conn)
	conn.channels[0].start()
	conn.handleOutgoing()
	conn.handleIncoming()
}

func (conn *Connection) hardClose() {
	conn.network.Close()
	conn.status.closed = true
	conn.server.deleteConnection(conn.id)
	conn.server.deleteRegisteredQueues(conn.id)
	for _, ch := range conn.channels {
		ch.shutdown()
	}
}

func (conn *Connection) closeConnWithError(err *proto.ProtoError) {
	fmt.Println("Sending connection close: ", err.Msg)
	conn.status.closing = true
	conn.channels[0].SendMethod(&proto.ConnectionClose{
		ReplyCode: err.Code,
		ReplyText: err.Msg,
		ClassId:   err.Class,
		MethodId:  err.Method,
	})
}

func (conn *Connection) removeChannel(chId uint16) {
	delete(conn.channels, chId)
}

func (conn *Connection) handleIncoming() {
	for {
		if conn.status.closed {
			break
		}
		frame, err := proto.ReadFrame(conn.network)
		if err != nil {
			fmt.Printf("Error reading frame: %s", err.Error())
			conn.hardClose()
			break
		}
		conn.handleFrame(frame)
	}
}

func (conn *Connection) handleOutgoing() {
	go func() {
		for {
			if conn.status.closed {
				break
			}
			frame := <-conn.outgoing
			proto.WriteFrame(frame)
		}
	}()
}

func (conn *Connection) handleFrame(f *proto.WireFrame) {
	// if !conn.status.open && f.Channel != 0 {
	if !conn.status.open {
		conn.hardClose()
		return
	}
	channel, ok := conn.channels[f.Channel]
	if !ok {
		channel = NewChannel(f.Channel, conn)
		conn.channels[f.Channel] = channel
		conn.channels[f.Channel].start()
	}
	// Dispatch frame to channel
	channel.incoming <- f
}
