package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sauravgsh16/message-server/allocate"
	"github.com/sauravgsh16/message-server/proto"
)

var counter int64

func init() {
	rand.Seed(time.Now().UnixNano())
	counter = time.Now().UnixNano()
}

func nextID() int64 {
	return atomic.AddInt64(&counter, 1)
}

// ConnectionStatus struct
type ConnectionStatus struct {
	start    bool
	startOk  bool
	open     bool
	openOk   bool
	closing  bool
	closed   bool
	closedOk bool
}

// Connection struct
type Connection struct {
	id        int64
	channels  map[uint16]*Channel
	outgoing  chan proto.Frame
	server    *Server
	network   net.Conn
	mux       sync.Mutex
	allocator allocate.Allocator
	status    ConnectionStatus
	writer    *proto.Writer
}

// NewConnection returns a new connection
func NewConnection(s *Server, n net.Conn) *Connection {
	return &Connection{
		id:       nextID(),
		channels: make(map[uint16]*Channel),
		outgoing: make(chan proto.Frame),
		server:   s,
		network:  n,
		status:   ConnectionStatus{},
		writer:   &proto.Writer{W: bufio.NewWriter(n)},
	}
}

func (c *Connection) openConnection() {

	// Protocol Handshake
	buf := make([]byte, 5)
	_, err := c.network.Read(buf)
	if err != nil {
		fmt.Printf("Error reading protocol header")
		c.hardClose()
		return
	}

	protoBytes := []byte{'S', 'E', 'C', 'O', 'C'}
	if bytes.Compare(buf, protoBytes) != 0 {
		// Write on connection, for client to send correct data for handshake
		c.network.Write(protoBytes)
		c.hardClose()
		return
	}
	// Create channel 0 and start the connection handshake

	c.channels[0] = NewChannel(0, c)
	c.channels[0].start()
	go c.handleIncoming(c.network)
	c.handleOutgoing()
}

func (c *Connection) hardClose() {
	if c.status.closed {
		return
	}

	c.network.Close()
	c.status.closed = true
	c.server.deleteConnection(c.id)
	c.server.deleteQueuesForConn(c.id)
	for _, ch := range c.channels {
		ch.shutdown()
	}
}

func (c *Connection) closeConnWithError(err *proto.Error) {
	fmt.Println("Sending connection close: ", err.Msg)
	c.status.closing = true
	c.channels[0].Send(&proto.ConnectionClose{
		ReplyCode: err.Code,
		ReplyText: err.Msg,
		ClassId:   err.Class,
		MethodId:  err.Method,
	})
}

func (c *Connection) removeChannel(chID uint16) {
	c.mux.Lock()
	delete(c.channels, chID)
	c.mux.Unlock()
}

func (c *Connection) send(f proto.Frame) error {
	if c.status.closed {
		return proto.NewHardError(500, "Sending on closed channel/Connection", 0, 0)
	}

	c.mux.Lock()
	err := c.writer.WriteFrame(f)
	c.mux.Unlock()
	if err != nil || c.status.closing == true {
		go c.hardClose()
	}
	return err
}

func (c *Connection) handleIncoming(r io.Reader) {

	buf := bufio.NewReader(r)
	frames := &proto.Reader{R: buf}

	for {
		if c.status.closed {
			break
		}
		frame, err := frames.ReadFrame()
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				continue
			}
			c.hardClose()
			break
		}
		c.handleFrame(frame)
	}
}

func (c *Connection) handleOutgoing() {
	go func() {
		for {
			if c.status.closed {
				break
			}
			frame := <-c.outgoing
			c.send(frame)
		}
	}()
}

func (c *Connection) handleFrame(f proto.Frame) {
	if !c.status.open && f.Channel() != 0 {
		c.hardClose()
		return
	}
	ch, ok := c.channels[f.Channel()]
	if !ok {
		ch = NewChannel(f.Channel(), c)
		c.channels[f.Channel()] = ch
		c.channels[f.Channel()].start()
	}
	// Dispatch frame to channel
	ch.incoming <- f
}
