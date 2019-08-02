package broker

import (
        "bufio"
        "io"
        "net"
        "errors"
        "sync"
        "time"
        "sync/atomic"
)

const (
        defaultConnTimeout = 30 * time.Second
)

var errMaxChannel = errors.New("max number of channels allocated")

type Connection struct {
        destructor sync.Once
        mux        sync.Mutex
        sendMux    sync.Mutex
        conn       io.ReadWriteCloser
        writer     *writer
        channels   map[int]*Channel
        allocator  *allocator
        closed     int32 // 1 for closed, 0 otherwise
        closes     []chan *Error
        errors     chan *Error
}

func (c *Connection) Channel() (*Channel, error) {
        id, ok := c.allocator.next()
        if !ok {
                return nil, errMaxChannel
        }
        ch := newChannel(c, uint16(id))
        c.channels[id] = ch
        return ch, nil
}

func (c *Connection) send(f frame) error {
        if c.IsClosed() {
                return ErrClosed
        }
        c.sendMux.Lock()
        err := c.writer.WriteFrame(f)
        c.sendMux.Unlock()
        if err != nil {
                go c.shutdown(&Error{
                        Code:   FrameError,
                        Reason: err.Error(),
                })
        }
        return err
}

func (c *Connection) shutdown(err *Error) {
        atomic.StoreInt32(&c.closed, 1)

        c.destructor.Do(func() {
                c.mux.Lock()
                defer c.mux.Unlock()

                if err != nil {
                        for _, c := range c.closes {
                                c <- err
                        }
                }

                if err != nil {
                        c.errors <- err
                }

                close(c.errors)

                for _, c := range c.closes {
                        close(c)
                }

                // Shutdown the channel
                for _, ch := range c.channels {
                        ch.shutdown(err)
                }
        })
} 

func (c *Connection) IsClosed() bool {
        return atomic.LoadInt32(&c.closed) == 1
}

func (c *Connection) demux(f frame) {
        // TODO
}

func (c *Connection) reader(r io.Reader) {
        buf := bufio.NewReader(r)
        frames := &reader{buf}

        for {
                frame, err := frames.ReadFrame()

                if err != nil {
                        c.shutdown(&Error{Code: FrameError, Reason: err.Error()})
                        return
                }
                c.demux(frame)
        }
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
        // Read data being written on the connection in separate goroutine
        go c.reader(conn)
        return c
}
