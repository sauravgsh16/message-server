package manager

import (
        "fmt"
        "bufio"
        "io"
)

type client struct {
        ID      int
        conn    io.ReadWriteCloser
        r       *bufio.Reader
        w       *bufio.Writer
        status  bool      // open - 0 or closed - 1
        closech chan interface{}
}

func newClient(id int, conn io.ReadWriteCloser) *client {
        return &client{
                ID:      id,
                r:       bufio.NewReader(conn),
                w:       bufio.NewWriter(conn),
                closech: make(chan interface{}),
                conn:    conn,
        }
}

func (c *client) Read(p []byte) (n int, err error) {
        //var delim byte = ','
        //b, err := c.r.ReadBytes(delim)
        n, err = c.r.Read(p)

        fmt.Printf("Here:  %v\n", p)
        /*
        if err != io.EOF {
                return 0, err
        }
        copy(p, b[:])
        */
        return
}

func (c *client) Write(p []byte) (n int, err error) {
        return c.w.Write(p)
}

func (c *client) Close() error {
        if err := c.conn.Close(); err != nil {
                return err
        }
        return nil
}