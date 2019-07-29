package exchange

import (
        "bufio"
        "io"
        "net"
        "errors"
        "sync"
        "time"
)

const (
        defaultConnTimeout = 30 * time.Second
)

var errMaxChannel = errors.New("max number of channels allocated")

type Connection struct {
        mux       sync.Mutex
        conn      io.ReadWriteCloser
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
        conn, err := dialer("tcp", addr, defaultConnTimeout)
        if err != nil {
                return nil, err
        }
        return open(conn), nil
}

func dialer(network, addr string, timeout time.Duration) (net.Conn, error) {
        conn, err := net.DialTimeout(network, addr, timeout)
        if err != nil {
                return nil, err
        }
        return conn, nil
}

func open(conn io.ReadWriteCloser) *Connection {
        c := &Connection{
                conn:      conn,
                channels:  make(map[int]*Channel),
                allocator: newAllocator(),
                writer:    &writer{bufio.NewWriter(conn)},
        }
        go c.reader(conn)
        return c
}

func (c *Connection) reader(r io.Reader) {
        buf := bufio.NewReader(r)
        _ = &reader{buf}
}
