package main

import (
	"fmt"
	"net"
	"os"

	"github.com/sauravgsh16/secoc-third/constant"
	"github.com/sauravgsh16/testtcp/server"
)

func handleConnection(sevr *server.Server, conn net.Conn) {
	sevr.OpenConnection(conn)
}

func main() {

	sevr := server.NewServer()
	ln, err := net.Listen("tcp", constant.UnsecuredPort)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Listening on port %s\n", constant.UnsecuredPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection\n")
			os.Exit(1)
		}
		fmt.Printf("Accepted conn: %+v\n", conn)
		go handleConnection(sevr, conn)
	}
}
