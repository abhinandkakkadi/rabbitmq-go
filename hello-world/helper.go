package main

import (
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
func connectToDB() *amqp.Connection {
	// TODO : read dsn from env
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	return conn

}
