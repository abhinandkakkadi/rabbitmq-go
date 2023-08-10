package main

import (
	"context"
	"log"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

// log errors
func failOnError(err error, msg string) {

	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}

}

// connect to rabbitMQ
func connectToMQ() *amqp.Connection {
	// TODO : read dsn from env
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn

}

func DeclareQueue(ch *amqp.Channel) amqp.Queue {

	q, err := ch.QueueDeclare(
		"task_queue", //name
		true,         // durable  -- even if rabbitMQ server goes down (or restart) - the message must be durable (not that if the consumer create the queue with same name first with durability false- then the queue will not be durable)
		false,        // delete when unused
		false,        //exclusive
		false,        //no-wait
		nil,          //arguments
	)
	failOnError(err, "Failed to declare a queue")

	return q

}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func PublishMessage(ctx context.Context, q amqp.Queue, ch *amqp.Channel, body string) {

	err := ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, //routing key
		false,  //mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")

}
