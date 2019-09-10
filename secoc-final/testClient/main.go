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
	_, err := qclient.Dial("tcp://localhost:9000")
	failOnError(err, "Failed to connect to qserver")

	fmt.Println("Success")
	// conn.Close()
}
