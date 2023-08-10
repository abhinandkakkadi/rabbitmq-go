package main

import (
	"context"
	"time"
)

func main() {

	conn := connectToDB()
	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", //name
		false,   // durable
		false,   // delete when unused
		false,   //exclusive
		false,   //no-wait
		nil,     //arguments
	)

	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
	defer cancel()

	body := "Hello World"
	err = ch.PublishWithContext(ctx,
		"",  // exchange
		q.Name, //routing key
		false,  //mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType : "text/plain",
			body
		}
	)


}
