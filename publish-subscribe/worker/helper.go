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


func DeclareExchange(ch *amqp.Channel) {

	err := ch.ExchangeDeclare(
		"logs",   // name  -- exchange name
		"fanout",  // type  -- exchange type
		true,  // durable
		false,  // auto-deleted
		false,  // internal
		false,  // no-wait
		nil,  // arguments
	)

	failOnError(err, "Failed to declare an exchange")

}

// Queue have to be declared here
// We have to make sure that the Queue exist in case if consumer starts before publisher
func DeclareQueue(ch *amqp.Channel) amqp.Queue {

	q, err := ch.QueueDeclare(
		"", //name -- nameless as 
		true,    // durable
		false,   // delete when unused
		false,   //exclusive
		false,   //no-wait
		nil,     //arguments
	)
	failOnError(err, "Failed to declare a queue")

	return q

}

func QueueBind(q amqp.Queue,ch *amqp.Channel) {

	err := ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
)
failOnError(err, "Failed to bind a queue")

}



func ConsumeMessage(ctx context.Context, q amqp.Queue, ch *amqp.Channel) {

	msgs, err := ch.Consume(
		q.Name, //queue
		"",     // consumer
		true,  // auto-acknowledgement -- if the worker dies before processing the work - the work (message) will be re queued
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	failOnError(err, "Failed to publish a message")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s",d.Body)
		}
	}()

	log.Printf(" [*] waiting for logs. To exit press CTRL+C")
	<- forever

}
