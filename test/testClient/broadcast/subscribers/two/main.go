package main

import (
	"fmt"
	"log"
	"time"

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

	err = ch.ExchangeDeclare(
		"test",   // name
		"fanout", // type
		false,    // noWait
	)

	q, err := ch.QueueDeclare(
		"qtest2", // name
		false,    // noWait
	)

	err = ch.QueueBind(
		q.Queue, // name
		"test",  // exhange name
		"",      // routing key
		false,   // noWait
	)

	msgs, err := ch.Consume(
		q.Queue, // Queue Name
		"c2",    // Consumer Name
		true,    // noAck
		false,   // noWait
	)
	failOnError(err, "Failed to register a consumer")

	var counter int
	timeout := time.After(30 * time.Second)

loop:
	for {
		select {
		case d := <-msgs:
			fmt.Printf("%+v\n", d)
			counter++
		case <-timeout:
			break loop
		}
	}
	fmt.Printf("Total messages received: %d\n", counter)
}
