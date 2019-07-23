package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/sauravgsh16/secoc-third/server/handler"
)

const (
	port = "9001"
)

func main() {
	if err := handler.Register(); err != nil {
		log.Fatalf("Failed to register handlers %v", err)
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Unable to start Queue server %v", err)
	} else {
		log.Printf("Queue server listening on port: %v", port)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Unable to accept connection on listener %v", err)
		}
		go rpc.ServeConn(conn)
	}
}