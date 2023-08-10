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

	DeclareExchange(ch)

	// call declare queue to declare a new queue
	q := DeclareQueue(ch)

	// bind exchange to queue in order for the exchange to forward the message to the queue
	QueueBind(q,ch)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ConsumeMessage(ctx, q, ch)

}
