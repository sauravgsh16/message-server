package main

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/sauravgsh16/message-server/constant"
	"github.com/sauravgsh16/message-server/qserver/server"
)

func handleConnection(sevr *server.Server, conn net.Conn) {
	sevr.OpenConnection(conn)
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get wd: %v", err)
		os.Exit(1)
	}
	serverDB := filepath.Join(wd, "server.db")
	msgStoreDB := filepath.Join(wd, "messages.db")

	sevr := server.NewServer(serverDB, msgStoreDB)
	ln, err := net.Listen("tcp", constant.UnsecuredPort)
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	log.Printf("Message server listening on port %s\n", constant.UnsecuredPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection")
			os.Exit(1)
		}
		log.Printf("Accepted conn: %+v\n", conn.LocalAddr())
		go handleConnection(sevr, conn)
	}
}
