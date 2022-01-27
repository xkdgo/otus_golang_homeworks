package main

import (
	"log"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/queue/controllers/rabbit"
)

const serviceNameExchange = "scheduler"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	consumer := rabbit.Consumer{}
	err := consumer.Init(
		rabbit.WithDialString("amqp://guest:guest@localhost:5672/"),
		rabbit.WithExchangeName(serviceNameExchange),
		rabbit.WithRoutingKey("calendar_sender"),
		rabbit.WithQueueName("sender"),
	)
	failOnError(err, "consumer init fail")
	msgs, err := consumer.ConsumeQueue()
	failOnError(err, "consumer fail to consume")

	consumererrorCh := consumer.GetErrorChannel()

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	errconsume := <-consumererrorCh
	log.Printf("%q\n", errconsume)
}
