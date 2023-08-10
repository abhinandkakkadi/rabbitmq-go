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


func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
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

func PublishMessage(ctx context.Context, ch *amqp.Channel, body string) {

	err := ch.PublishWithContext(ctx,
		"logs",     // exchange
		"", //routing key
		false,  //mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")

}
