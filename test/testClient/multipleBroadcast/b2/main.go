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

	done := make(chan interface{})

	err = ch.ExchangeDeclare(
		"test",   // name
		"fanout", // type
		false,    // noWait
	)
	failOnError(err, "Failed to declare exchange")

	after := time.After(10 * time.Second)
	body := []byte("This is a test string 2")

	go func() {
		select {
		case <-after:
			fmt.Println("Publishing 2")
			err = ch.Publish(
				"test", // name
				"",     // routing key
				false,  // immediate
				qclient.MetaDataWithBody{
					ContentType:   "text/plain",
					MessageID:     "msgid123",
					UserID:        "userid123",
					ApplicationID: "appid123",
					Body:          body,
				},
			)
			failOnError(err, "Failed to publish")
			done <- true
		}
	}()

	<-done
	fmt.Printf("2st publish")
}
