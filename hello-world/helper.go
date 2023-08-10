package main

import (
	"context"
	"log"

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
		"hello", //name
		false,   // durable
		false,   // delete when unused
		false,   //exclusive
		false,   //no-wait
		nil,     //arguments
	)
	failOnError(err, "Failed to declare a queue")

	return q

}

func PublishMessage(ctx context.Context,q amqp.Queue,ch *amqp.Channel, body string) {
	
	
	err := ch.PublishWithContext(ctx,
		"",  // exchange
		q.Name, //routing key
		false,  //mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType : "text/plain",
			Body: []byte(body),
		})
		failOnError(err,"Failed to publish a message")
		
}
