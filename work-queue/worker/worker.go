package main

import (
	"context"
	"time"
)

func main() {

	// connect tp service
	conn := connectToMQ()
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")

	defer ch.Close()

	// call declare queue to declare a new queue
	q := DeclareQueue(ch)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ConsumeMessage(ctx, q, ch)

}
