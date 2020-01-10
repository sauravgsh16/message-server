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

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // no-wait
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World! from 1"
	after := time.After(5 * time.Second)
	done := make(chan interface{})

	go func() {
		select {
		case <-after:
			log.Println("Publishing in 1")
			err = ch.Publish(
				"",      // exchange
				q.Queue, // routing key
				false,   // immediate
				qclient.MetaDataWithBody{
					ContentType:   "text/plain",
					MessageID:     "msgid123",
					UserID:        "userid123",
					ApplicationID: "appid123",
					Body:          []byte(body),
				},
			)
			failOnError(err, "Failed to publish a message")
			close(done)
		}
	}()

	<-done
	fmt.Println("Successfully sent")
}
