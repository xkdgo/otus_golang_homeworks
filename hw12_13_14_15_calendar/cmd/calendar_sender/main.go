package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// err = ch.ExchangeDeclare(
	// 	"logs_direct", // name
	// 	"direct",      // type
	// 	true,          // durable
	// 	false,         // auto-deleted
	// 	false,         // internal
	// 	false,         // no-wait
	// 	nil,           // arguments
	// )
	// failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"sender", // name
		false,    // durable
		false,    // delete when unused
		true,     // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_direct", "calendar_sender")
	err = ch.QueueBind(
		q.Name,            // queue name
		"calendar_sender", // routing key
		"scheduler",       // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name,            // queue
		"calendar_sender", // consumer
		true,              // auto ack
		false,             // exclusive
		false,             // no local
		false,             // no wait
		nil,               // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
