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

// Queue have to be declared here
// We have to make sure that the Queue exist in case if consumer starts before publisher
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

func ConsumeMessage(ctx context.Context,q amqp.Queue,ch *amqp.Channel) {
	
	
		msgs, err := ch.Consume(
			q.Name, //queue
			"", // consumer
			true,  // auto-acknowledgement
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,  // args
		)

		failOnError(err,"Failed to publish a message")

		var forever chan struct{}

		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
			}
		}()

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<- forever
		
}
