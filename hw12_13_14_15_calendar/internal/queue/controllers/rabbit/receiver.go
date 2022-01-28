package rabbit

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/queue"
)

type Receiver struct {
	opts         queue.Options
	dialString   string
	exchangeName string
	queueName    string
	routingKey   string
	connection   *amqp.Connection
	channel      *amqp.Channel
	wg           *sync.WaitGroup
	log          *logger.Logger
	handler      Worker
}

func (c *Receiver) Init(opts ...queue.Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	if c.opts.Logger != nil {
		c.log = c.opts.Logger
	} else {
		return fmt.Errorf("init fail missed Logger")
	}
	if handler, ok := c.opts.Context.Value(handlerFunKey{}).(Worker); ok {
		c.handler = handler
	} else {
		return fmt.Errorf("init fail missed exchangeName")
	}
	if exchangeName, ok := c.opts.Context.Value(exchangeNameKey{}).(string); ok {
		c.exchangeName = exchangeName
	} else {
		return fmt.Errorf("init fail missed exchangeName")
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

func (c *Receiver) Dial() error {
	c.log.Debugf("dial to %s", c.dialString)
	conn, err := amqp.Dial(c.dialString)
	if err != nil {
		return err
	}
	c.connection = conn
	return nil
}

func (c *Receiver) ConsumeQueue() (deliv <-chan amqp.Delivery, err error) {
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

func (c *Receiver) GetErrorChannel() chan *amqp.Error {
	errCh := make(chan *amqp.Error)
	errCh = c.channel.NotifyClose(errCh)
	return errCh
}

func (c *Receiver) reconnect(ctx context.Context, timeout time.Duration) {
	err := c.Dial()
	for err != nil {
		ticker := time.NewTicker(time.Second)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err = c.Dial()
			if err != nil {
				c.log.Infof(" [err] failed reconnect %v", err)
			} else {
				c.log.Infof("connected ...")
			}

		}
	}
}

func (c *Receiver) Handle(
	ctx context.Context,
	numWorkers int,
	reconnectTimeout time.Duration) {
	needToExit := make(chan struct{}, 1)
	go func() {
		<-ctx.Done()
		close(needToExit)
	}()
	for {
		select {
		case <-needToExit:
			return
		default:
			if c.connection.IsClosed() {
				c.reconnect(ctx, reconnectTimeout)
			}
			msgs, err := c.ConsumeQueue()
			for err != nil && !c.connection.IsClosed() {
				ticker := time.NewTicker(reconnectTimeout)
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					c.log.Infof(" [err] %q", err)
					msgs, err = c.ConsumeQueue()
				}
			}
			consumererrorCh := c.GetErrorChannel()
			closeCh := make(chan struct{}, 1)
			for i := 0; i < numWorkers; i++ {
				c.wg.Add(1)
				go func(ctx context.Context, i int, closeCh chan struct{}) {
					c.log.Debugf("gouroutine %d started", i)
					defer c.wg.Done()
					defer c.log.Debugf("gouroutine %d exited", i)
					for {
						select {
						case <-ctx.Done():
							c.log.Debugf("gouroutine %d context closed", i)
							return
						case <-closeCh:
							c.log.Debugf("gouroutine %d catched closeCh event", i)
							return
						case d, open := <-msgs:
							if !open {
								return
							}
							c.handler(ctx, d.ContentEncoding, d.Body)
						}
					}
				}(ctx, i, closeCh)
			}
			go func() {
				errconsume := <-consumererrorCh
				if errconsume == nil {
					c.log.Infof("listen catch nil error %v and exit", errconsume)
					needToExit <- struct{}{}
					close(closeCh)
					return
				}
				c.log.Infof("%v\n", errconsume)
				close(closeCh)
			}()
			c.log.Infof(" [*] Waiting for msgs. To exit press CTRL+C")
			c.wg.Wait()
			time.Sleep(time.Second)
		}
	}
}
