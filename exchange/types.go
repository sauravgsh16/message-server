package exchange

import (
        "io"
)

type message interface {
      read(io.Reader) error
      write(io.Writer) error  
}

type writer struct {
        w io.Writer
}

type reader struct {
        r io.Reader
}

type exchangeDeclare struct {
        Exchange string
        Type     string
}

type exchangeDeclareOk struct {}
