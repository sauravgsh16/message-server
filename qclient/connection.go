package qclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/sauravgsh16/secoc-third/allocate"
	"github.com/sauravgsh16/secoc-third/proto"
)

const defaultConnTimeout = 30 * time.Second

var (
	ErrMaxChannel     = errors.New("max number of channels allocated")
	ErrClosed         = errors.New("Send attempt on closed connection")
	ErrInvalidCommand = errors.New("Invalid command received")
	ErrGetHostName    = errors.New("Unable to retrieve Host Name")
	ErrHost           = errors.New("Corrupt HostName")
)

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
	destructor sync.Once
	sendMux    sync.Mutex
	mux        sync.Mutex
	conn       io.ReadWriteCloser
	channels   map[int]*Channel
	outgoing   chan *proto.WireFrame
	incoming   chan *proto.MethodFrame
	status     ConnectionStatus
	statusMux  sync.RWMutex
	errors     chan *proto.Error
	allocator  *allocate.Allocator
	version    byte
	mechanisms string
}

func Dial(url string) (*Connection, error) {
	return dial(url)
}

func dial(url string) (*Connection, error) {
	uri, err := parseUrl(url)
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

func Open(conn io.ReadWriteCloser) (*Connection, error) {
	c := &Connection{
		conn:      conn,
		channels:  make(map[int]*Channel),
		outgoing:  make(chan *proto.WireFrame, 100),
		incoming:  make(chan *proto.MethodFrame),
		allocator: allocate.NewAllocator(),
		status:    ConnectionStatus{},
	}
	c.handleOutgoing()
	go c.handleIncoming()
	return c, c.open()
}

// Connection method receivers
func (c *Connection) IsClosed() bool {
	c.statusMux.Lock()
	defer c.statusMux.Unlock()
	return c.status.closed
}

func (c *Connection) send(cf *proto.ChannelFrame) error {
	if c.IsClosed() {
		return ErrClosed
	}

	buf := bytes.NewBuffer([]byte{})
	c.sendMux.Lock()
	err := cf.Method.Write(buf)
	c.sendMux.Unlock()

	if err != nil {
		go c.hardClose()
	}
	c.outgoing <- &proto.WireFrame{
		FrameType: uint8(proto.FrameMethod),
		Channel:   cf.ChannelID,
		Payload:   buf.Bytes(),
	}
	return err
}

func (c *Connection) call(req proto.MethodFrame, res ...proto.MethodFrame) error {
	if req != nil {
		if err := c.send(&proto.ChannelFrame{ChannelID: uint16(0), Method: req}); err != nil {
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
		for _, result := range res {
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
	return nil
}

func (c *Connection) open() error {
	header := &proto.ChannelFrame{
		ChannelID: uint16(0),
		Method:    &proto.ProtocolHeader{},
	}
	if err := c.send(header); err != nil {
		return err
	}
	return c.openStart()
}

func (c *Connection) openStart() error {
	start := &proto.ConnectionStart{}

	if err := c.call(nil, start); err != nil {
		return err
	}
	c.version = start.Version
	c.mechanisms = start.Mechanisms

	startOk := &proto.ConnectionStartOk{
		Mechanism: c.mechanisms,
		Response:  "Authenticated",
	}
	if err := c.send(&proto.ChannelFrame{ChannelID: uint16(0), Method: startOk}); err != nil {
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
	return nil
}

func (c *Connection) handleIncoming() {
	for {
		if c.status.closed {
			break
		}
		frame, err := proto.ReadFrame(c.conn)
		if err != nil {
			fmt.Printf("Error reading frame: %s", err.Error())
			c.hardClose()
			break
		}
		c.handleFrame(frame)
	}
}

func (c *Connection) handleOutgoing() {
	go func() {
		for {
			if c.IsClosed() {
				break
			}
			frame := <-c.outgoing
			proto.WriteFrame(c.conn, frame)
		}
	}()
}

func (c *Connection) Channel() (*Channel, error) {
	id, ok := c.allocator.Next()
	if !ok {
		return nil, ErrMaxChannel
	}
	ch := newChannel(c, uint16(id))
	c.channels[id] = ch
	return ch, nil
}
