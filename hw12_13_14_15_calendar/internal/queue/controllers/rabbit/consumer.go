package rabbit

import (
	"context"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/queue"
)

type Consumer struct {
	opts         queue.Options
	dialString   string
	exchangeName string
	queueName    string
	routingKey   string
	connection   *amqp.Connection
	channel      *amqp.Channel
	wg           *sync.WaitGroup
}

func (c *Consumer) Init(opts ...queue.Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	if exchangeName, ok := c.opts.Context.Value(exchangeNameKey{}).(string); ok {
		c.exchangeName = exchangeName
	} else {
		return fmt.Errorf("init fail missed exchangeName")
	}
	if queueName, ok := c.opts.Context.Value(queueNameKey{}).(string); ok {
		c.queueName = queueName
	} else {
		return fmt.Errorf("init fail missed queueName")
	}
	if routingKey, ok := c.opts.Context.Value(routingKeyKey{}).(string); ok {
		c.routingKey = routingKey
	} else {
		return fmt.Errorf("init fail missed routingKey")
	}
	if dialString, ok := c.opts.Context.Value(dialStringKey{}).(string); ok {
		c.dialString = dialString
	} else {
		return fmt.Errorf("init fail missed dialString")
	}
	err := c.Dial()
	if err != nil {
		return err
	}
	c.wg = &sync.WaitGroup{}
	return nil
}

func (c *Consumer) Dial() error {
	conn, err := amqp.Dial(c.dialString)
	if err != nil {
		return err
	}
	c.connection = conn
	return nil
}

func (c *Consumer) ConsumeQueue() (deliv <-chan amqp.Delivery, err error) {
	ch, err := c.connection.Channel()
	if err != nil {
		return nil, err
	}
	c.channel = ch
	_, err = c.channel.QueueDeclare(
		c.queueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, err
	}
	err = c.channel.QueueBind(
		c.queueName,    // queue name
		c.routingKey,   // routing key
		c.exchangeName, // exchange
		false,
		nil)
	if err != nil {
		return nil, err
	}
	deliv, err = c.channel.Consume(
		c.queueName,  // queue
		c.routingKey, // consumer
		true,         // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	if err != nil {
		return nil, err
	}
	return deliv, nil
}

func (c *Consumer) GetErrorChannel() chan *amqp.Error {
	errCh := make(chan *amqp.Error)
	errCh = c.channel.NotifyClose(errCh)
	return errCh
}

type Worker func(ctx context.Context, msg queue.NotifyEvent, wg *sync.WaitGroup)

// func (c *Consumer) Handle(ctx context.Context, fn Worker, numWorkers int) {
// 	msgs, err := c.ConsumeQueue()
// 	failOnError(err, "consumer fail to consume")

// 	consumererrorCh := consumer.GetErrorChannel()

// 	go func() {
// 		for d := range msgs {
// 			log.Printf(" [x] %s", d.Body)
// 		}
// 	}()

// 	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
// 	errconsume := <-consumererrorCh
// 	log.Printf("%q\n", errconsume)
// }
