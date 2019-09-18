package main

import (
	"fmt"
	"log"

	"github.com/sauravgsh16/secoc-third/qclient"
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

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	err = ch.ExchangeDeclare(
		"test",   // name
		"fanout", // type
		false,    // noWait
	)
	failOnError(err, "Failed to declare exchange")

	body := []byte("This is a test string")

	err = ch.Publish(
		"test", // name
		"",     // routing key
		false,  // immediate
		body,
	)
	failOnError(err, "Failed to publish a message")

	fmt.Println("Success")
}
