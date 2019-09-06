package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
)

type Connection struct {
	server *Server
	conn   io.ReadWriteCloser
	rpc    chan message
	writer *writer
}

func NewConn(s *Server, n net.Conn) *Connection {
	return &Connection{
		server: s,
		conn:   n,
		rpc:    make(chan message),
		writer: &writer{bufio.NewWriter(n)},
	}
}

func (c *Connection) openConn() {
	go c.handleIncoming(c.conn)
	c.start()
}

func (c *Connection) send(f frame) error {
	return c.writer.WriteFrame(f)
}

func (c *Connection) call(req message, resp ...message) error {
	if req != nil {
		if err := c.send(&methodFrame{ChannelId: 0, Method: req}); err != nil {
			return err
		}
	}

	select {
	case msg := <-c.rpc:
		// Try to match one of the result types
		for _, try := range resp {
			if reflect.TypeOf(msg) == reflect.TypeOf(try) {
				// *res = *msg
				vres := reflect.ValueOf(try).Elem()
				vmsg := reflect.ValueOf(msg).Elem()
				vres.Set(vmsg)
				return nil
			}
		}
		return errors.New("ERRORRRRRRR!!!!!!")
	}
}

func (c *Connection) start() {
	start := &ConnectionStart{
		Version:    1,
		Mechanisms: "PLAIN",
	}
	ok := &ConnectionStartOk{}

	c.call(start, ok)

	fmt.Printf("%+v", ok)
}

func (c *Connection) demux(f frame) {
	switch mf := f.(type) {
	case *methodFrame:
		switch m := mf.Method.(type) {
		default:
			c.rpc <- m
		}
	}
}

func (c *Connection) handleIncoming(r io.Reader) {
	buf := bufio.NewReader(r)
	frames := &reader{r: buf}

	for {
		frame, err := frames.ReadFrame()

		if err != nil && err != io.EOF {
			fmt.Println("Error reading frame" + err.Error())
		}

		c.demux(frame)
		err = nil
	}
}
