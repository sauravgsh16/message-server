package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"time"
)

type Connection struct {
	conn   io.ReadWriteCloser
	rpc    chan message
	writer *writer
}

func Dial(url string) (*Connection, error) {
	conn, err := net.DialTimeout("tcp", "localhost:9000", 30*time.Second)
	if err != nil {
		return nil, err
	}
	return open(conn)
}

func open(conn io.ReadWriteCloser) (*Connection, error) {
	c := &Connection{
		conn:   conn,
		rpc:    make(chan message),
		writer: &writer{bufio.NewWriter(conn)},
	}
	go c.handleIncoming(conn)
	return c, c.openstart()
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

func (c *Connection) send(f frame) error {
	return c.writer.WriteFrame(f)
}

func (c *Connection) openstart() error {
	start := &ConnectionStart{}

	if err := c.call(nil, start); err != nil {
		return err
	}

	ok := &ConnectionStartOk{
		Mechanism: "PLAIN",
		Response:  "AUTH",
	}

	fmt.Printf("%+v", start)

	return c.send(&methodFrame{ChannelId: 0, Method: ok})
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

	count := 0
	for {
		frame, err := frames.ReadFrame()

		if err != nil && err != io.EOF {
			fmt.Printf("Error reading frame: %s\n", err.Error())
		}

		c.demux(frame)
		count++
		fmt.Printf("Count: %d\n", count)
	}
}

func (c *Connection) Close() {
	c.conn.Close()
}
