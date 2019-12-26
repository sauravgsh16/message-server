package main

import (
	"log"

	"github.com/sauravgsh16/message-server/qclient"
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

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // no-wait
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Queue, // queue
		"c1",    // consumer ** Need to pass **
		true,    // noAck
		false,   // noWait
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
}
