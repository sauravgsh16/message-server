package exchange

import (
        "net"
        "net/rpc"
        "errors"
)

var errMaxChannel = errors.New("max number of channels allocated")

type Connection struct {
        conn      *rpc.Client
        writer    *writer
        channels  map[int]*Channel
        allocator *allocator
}

func (c *Connection) Channel() (*Channel, error) {
        id, ok := c.allocator.next()
        if !ok {
                return nil, errMaxChannel
        }
        ch := newChannel(c, id)
        c.channels[id] = ch
        return ch, nil
}

func (c *Connection) Close() {
        // TODO
}

// Dail connects exchange to server
func Dial(url string) (*Connection, error) {
        return dial(url)
}

func dial(url string) (*Connection, error) {
        u, err := parseUri(url)
        if err != nil {
                return nil, err
        }
        addr := net.JoinHostPort(u.host, u.port)
        conn, err := dialer("tcp", addr)
        if err != nil {
                return nil, err
        }
        return open(conn), nil
}

func dialer(network, addr string) (*rpc.Client, error) {
        conn, err := rpc.Dial(network, addr)
        if err != nil {
                return nil, err
        }
        return conn, nil
}

func open(conn *rpc.Client) *Connection {
        c := &Connection{
                conn:      conn,
                channels:  make(map[int]*Channel),
                allocator: newAllocator(),
        }
        return c
}
