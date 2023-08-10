package main

import (
	"bytes"
	"context"
	"log"
	"time"

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
		true,    // durable
		false,   // delete when unused
		false,   //exclusive
		false,   //no-wait
		nil,     //arguments
	)
	failOnError(err, "Failed to declare a queue")

	return q

}

func ConsumeMessage(ctx context.Context, q amqp.Queue, ch *amqp.Channel) {

	msgs, err := ch.Consume(
		q.Name, //queue
		"",     // consumer
		false,  // auto-acknowledgement -- if the worker dies before processing the work - the work (message) will be re queued
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	failOnError(err, "Failed to publish a message")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false) // acknowledge a single delivery - once we're done with a task
			// Using this code, you can ensure that even if you terminate a worker using CTRL+C while it was processing a message, nothing is lost. Soon after the worker terminates, all unacknowledged messages are redelivered.
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C ")
	<-forever

}
