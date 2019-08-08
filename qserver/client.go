package qserver

import (
        "bufio"
        "io"
)

type Client struct {
        ID      int
        r       io.Reader
        w       io.Writer
        status  bool      // open - 0 or closed - 1
        closech chan interface{}
}

func NewClient(id int, conn io.ReadWriter) *Client {
        return &Client{
                ID:      id,
                r:       bufio.NewReader(conn),
                w:       bufio.NewWriter(conn),
                closech: make(chan interface{}),
        }
}

func (c *Client) Read(p []byte) (n int, err error) {
        return c.r.Read(p)
}

func (c *Client) Write(p []byte) (n int, err error) {
        return c.w.Write(p)
}