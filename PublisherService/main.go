package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {

	sendSignal("Hello World")
}

func sendSignal(body string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()

	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	FailOnError(err, "Failed to publish a message")
}

func FailOnError(err error, msg string) {
	if err != nil {
		println(msg + ": " + err.Error())
	}
}
