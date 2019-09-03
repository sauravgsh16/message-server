package server

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sauravgsh16/secoc-third/allocate"
	"github.com/sauravgsh16/secoc-third/proto"
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

	// Protocol Handshake
	buf := make([]byte, 5)
	_, err := conn.network.Read(buf)
	if err != nil {
		fmt.Printf("Error reading protocol header")
		conn.hardClose()
		return
	}

	protoBytes := []byte{'S', 'E', 'C', 'O', 'C'}
	if bytes.Compare(buf, protoBytes) != 0 {
		// Write on connection, for client to send correct data for handshake
		conn.network.Write(protoBytes)
		conn.hardClose()
		return
	}

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
	conn.server.deleteQueuesForConn(conn.id)
	for _, ch := range conn.channels {
		ch.shutdown()
	}
}

func (conn *Connection) closeConnWithError(err *proto.Error) {
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
			proto.WriteFrame(conn.network, frame)
		}
	}()
}

func (conn *Connection) handleFrame(wf *proto.WireFrame) {
	// if !conn.status.open && f.Channel != 0 {
	if !conn.status.open {
		conn.hardClose()
		return
	}
	ch, ok := conn.channels[wf.Channel]
	if !ok {
		ch = NewChannel(wf.Channel, conn)
		conn.channels[wf.Channel] = ch
		conn.channels[wf.Channel].start()
	}
	// Dispatch frame to channel
	ch.incoming <- wf
}
