package main

import (
	"fmt"
	"log"

	"github.com/sauravgsh16/secoc-third/secoc-final/qclient"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := qclient.Dial("tcp://localhost:9000")
	// defer conn.Close()
	failOnError(err, "Failed to connect to qserver")

	_, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	//defer ch.Close()

	fmt.Println("Success")
}
