package qutils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const SensorListQueue = "SensorList"

func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return conn, ch
}

func GetQueue(name string, ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(
		name,  // q name
		false, // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	failOnError(err, "Failed to declare a queue")
	return &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
