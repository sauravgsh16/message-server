package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("localhost:9000")
	failOnError(err, "Failed to connect to qserver")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	/*
			err = ch.ExchangeDeclare(
				"logs",   // name
				"fanout", // type
				true,     // durable
				false,    // auto-deleted
				false,    // internal
				false,    // no-wait
				nil,      // arguments
			)
			failOnError(err, "Failed to declare an exchange")

			body := bodyFrom(os.Args)
			err = ch.Publish(
				"logs", // exchange
				"",     // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			failOnError(err, "Failed to publish a message")

		        log.Printf(" [x] Sent %s", body)
	*/
}
