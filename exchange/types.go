package exchange

import (
        "io"
)

type message interface {
      read(io.Reader) error
      write(io.Writer) error  
}

type writer struct {
        io.Writer
}

type reader struct {
        io.Reader
}