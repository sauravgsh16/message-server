package qclient

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/sauravgsh16/message-server/allocate"
	"github.com/sauravgsh16/message-server/proto"
)

const defaultConnTimeout = 30 * time.Second

var (
	ErrMaxChannel     = errors.New("max number of channels allocated")
	ErrInvalidCommand = errors.New("Invalid command received")
	ErrGetHostName    = errors.New("Unable to retrieve Host Name")
	ErrHost           = errors.New("Corrupt HostName")
)

// ConnectionStatus represents connection status
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
	destructor sync.Once
	sendMux    sync.Mutex
	mux        sync.Mutex
	conn       io.ReadWriteCloser
	channels   map[uint16]*Channel
	outgoing   chan proto.Frame
	incoming   chan proto.MessageFrame
	status     ConnectionStatus
	statusMux  sync.RWMutex
	errors     chan *proto.Error
	allocator  *allocate.Allocator
	writer     *proto.Writer
}

// Dial to connect to a listener
func Dial(url string) (*Connection, error) {
	return dial(url)
}

func dial(url string) (*Connection, error) {
	uri, err := parseURL(url)
	if err != nil {
		return nil, err
	}
	addr := net.JoinHostPort(uri.host, uri.port)
	conn, err := dialer("tcp", addr, defaultConnTimeout)
	if err != nil {
		return nil, err
	}
	return Open(conn)
}

func dialer(netType, addr string, timeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout(netType, addr, timeout)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Open a connection
func Open(conn io.ReadWriteCloser) (*Connection, error) {
	c := &Connection{
		conn:     conn,
		channels: make(map[uint16]*Channel),
		outgoing: make(chan proto.Frame),
		incoming: make(chan proto.MessageFrame),
		errors:   make(chan *proto.Error, 1),
		status:   ConnectionStatus{},
		writer:   &proto.Writer{W: bufio.NewWriter(conn)},
	}
	go c.handleOutgoing()
	go c.handleIncoming(c.conn)
	return c, c.open()
}

// IsClosed return if connection is closed
func (c *Connection) IsClosed() bool {
	c.statusMux.Lock()
	defer c.statusMux.Unlock()
	return c.status.closed
}

// Close connection
func (c *Connection) Close() error {
	if c.IsClosed() {
		return ErrClosed
	}

	err := c.call(
		&proto.ConnectionClose{
			ReplyCode: 200,
			ReplyText: "Bye",
		},
		&proto.ConnectionCloseOk{},
	)
	c.hardClose(nil)
	return err
}

func (c *Connection) send(f proto.Frame) error {
	if c.status.closed {
		return proto.NewHardError(500, "Sending on closed channel/Connection", 0, 0)
	}
	c.mux.Lock()
	err := c.writer.WriteFrame(f)
	c.mux.Unlock()
	if err != nil {
		pErr := proto.NewHardError(500, err.Error(), 0, 0)
		go c.hardClose(pErr)
	}
	return err
}

func (c *Connection) call(req proto.MessageFrame, resp ...proto.MessageFrame) error {
	if req != nil {

		fmt.Println("Sending", req.MethodName())

		if err := c.send(&proto.MethodFrame{ChannelID: uint16(0), Method: req}); err != nil {
			return err
		}
	}

	select {
	case err, ok := <-c.errors:
		if !ok {
			return ErrClosed
		}
		return err
	case msg := <-c.incoming:
		// We try to match the 'res' - result types
		for _, result := range resp {

			if reflect.TypeOf(msg) == reflect.TypeOf(result) {
				// we updates res with the data in result
				// Thus making, *result = *msg
				vres := reflect.ValueOf(result).Elem()
				vmsg := reflect.ValueOf(msg).Elem()
				vres.Set(vmsg)

				return nil
			}
		}

		return ErrInvalidCommand
	}
}

func (c *Connection) open() error {
	if err := c.send(&proto.ProtocolHeader{}); err != nil {
		return err
	}

	return c.openStart()
}

func (c *Connection) openStart() error {
	start := &proto.ConnectionStart{}

	if err := c.call(nil, start); err != nil {
		return err
	}

	startOk := &proto.ConnectionStartOk{
		Mechanism: start.Mechanisms,
		Response:  "Authenticated",
	}
	if err := c.send(&proto.MethodFrame{ChannelID: uint16(0), Method: startOk}); err != nil {
		return err
	}
	return c.openHost()
}

func (c *Connection) openHost() error {
	host, err := os.Hostname()
	if err != nil {
		return ErrGetHostName
	}
	req := &proto.ConnectionOpen{Host: host}
	res := &proto.ConnectionOpenOk{}

	if err := c.call(req, res); err != nil {
		return ErrHost
	}
	c.allocator = allocate.NewAllocator()
	return nil
}

func (c *Connection) hardClose(err *proto.Error) {
	c.status.closing = true

	c.destructor.Do(func() {
		c.mux.Lock()
		defer c.mux.Unlock()

		if err != nil {
			c.errors <- err
		}
		close(c.errors)

		for _, ch := range c.channels {
			ch.shutdown(err)
		}

		c.conn.Close()

		c.channels = map[uint16]*Channel{}
		c.allocator = allocate.NewAllocator()
	})
}

func (c *Connection) closeWithErr(err *proto.Error) {
	if c.IsClosed() {
		return
	}
	defer c.hardClose(err)

	c.call(&proto.ConnectionClose{
		ReplyCode: err.Code,
		ReplyText: err.Msg,
		ClassId:   err.Class,
		MethodId:  err.Method,
	}, &proto.ConnectionCloseOk{})
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
			pErr := proto.NewHardError(500, err.Error(), 0, 0)
			c.hardClose(pErr)
			break
		}
		if frame != nil {
			c.handleFrame(frame)
		}
	}
}

func (c *Connection) handleOutgoing() {
	for {
		if c.status.closed {
			break
		}
		frame := <-c.outgoing
		c.send(frame)
	}
}

func (c *Connection) handleFrame(f proto.Frame) {
	if f.Channel() == 0 {
		c.dispatch0(f)
	} else {
		c.dispatchN(f)
	}
}

func (c *Connection) dispatch0(f proto.Frame) {
	switch mf := f.(type) {

	case *proto.MethodFrame:
		c.routeMethod(mf)

	default:
		c.closeWithErr(ErrUnexpectedFrame)
	}
}

func (c *Connection) dispatchN(f proto.Frame) {
	c.mux.Lock()
	ch := c.channels[f.Channel()]
	c.mux.Unlock()

	if ch != nil {
		// Send data to channel to be processed
		ch.incoming <- f
	} else {
		// We expect the method here to be ConnectionCloseOk, ChannelClose, or ChannelCloseOk
		c.routeMethod(f.(*proto.MethodFrame))
	}
}

func (c *Connection) routeMethod(mf *proto.MethodFrame) *proto.Error {

	fmt.Println("Received", mf.Method.MethodName())

	clsID, mtdID := mf.Method.Identifier()

	switch clsID {
	case 10:
		switch method := mf.Method.(type) {

		case *proto.ConnectionClose:
			c.send(&proto.MethodFrame{
				ChannelID: uint16(0),
				Method:    &proto.ConnectionCloseOk{},
			})
		default:
			c.incoming <- method
		}
	case 20:
		switch mf.Method.(type) {
		case *proto.ChannelClose:
			c.send(&proto.MethodFrame{
				ChannelID: uint16(mf.ChannelID),
				Method:    &proto.ChannelCloseOk{},
			})
		case *proto.ChannelCloseOk:
			// Case ChannelCloseOk, since channel already closed, we ignore.
		default:
			// Unexpected method
			err := proto.NewHardError(504, "Communication attempt on close Channel/Connection", clsID, mtdID)
			c.closeWithErr(err)
		}
	}
	return nil
}

func (c *Connection) allocateChannel() (*Channel, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.IsClosed() {
		return nil, ErrClosed
	}

	id, ok := c.allocator.Next()
	if !ok {
		return nil, ErrMaxChannel
	}

	ch := newChannel(c, uint16(id))
	c.channels[uint16(id)] = ch
	return ch, nil
}

func (c *Connection) releaseChannel(id uint16) {
	c.mux.Lock()
	defer c.mux.Unlock()

	delete(c.channels, id)
}

func (c *Connection) openChannel() (*Channel, error) {
	ch, err := c.allocateChannel()
	if err != nil {
		return nil, err
	}

	if err := ch.open(); err != nil {
		c.releaseChannel(ch.id)
		return nil, err
	}
	return ch, nil
}

func (c *Connection) closeChannel(ch *Channel, err *proto.Error) {
	ch.shutdown(err)
	c.releaseChannel(ch.id)
}

// Channel opens a channel for the connection
func (c *Connection) Channel() (*Channel, error) {
	return c.openChannel()
}
