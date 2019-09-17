package main

import (
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

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	err = ch.ExchangeDeclare(
		"test",   // name
		"fanout", // type
		false,    // noWait
	)

	q, err := ch.QueueDeclare(
		"qtest", // name
		false,   // noWait
	)

	err = ch.QueueBind(
		q.Queue, // name
		"test",  // exhange name
		"",      // routing key
		false,   // noWait
	)

}
