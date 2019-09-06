package main

import (
	"fmt"
	"log"

	"github.com/sauravgsh16/testtcp/client"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := client.Dial("tcp://localhost:9000")
	failOnError(err, "Failed to connect to qserver")

	fmt.Println("Success")
	conn.Close()
}
