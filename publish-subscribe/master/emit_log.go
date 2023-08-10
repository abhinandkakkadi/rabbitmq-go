package main

import (
	"context"
	"log"
	"os"
	"time"
)

func main() {

	// connect tp service
	conn := connectToMQ()
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")

	defer ch.Close()


	DeclareExchange(ch)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// p
	body := bodyFrom(os.Args)

	PublishMessage(ctx,ch, body)

	// push/publish a message to queue
	log.Printf(" [x] Sent %s\n", body)
}
