package app

import (
	"context"
	"log"
	"time"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/queue"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/queue/controllers/rabbit"
)

type Sender struct {
	logger           *logger.Logger
	numWorkers       int
	dialString       string
	exchangeName     string
	routingKey       string
	queueName        string
	reconnectTimeOut time.Duration
	Subscriber       queue.Subscriber
}

func (a *Sender) Init() error {
	err := a.Subscriber.Init(
		queue.WithLogger(a.logger),
		rabbit.WithHandlerFunc(a.handler),
		rabbit.WithDialString(a.dialString),
		rabbit.WithExchangeName(a.exchangeName),
		rabbit.WithRoutingKey(a.routingKey),
		rabbit.WithQueueName(a.queueName),
	)
	if err != nil {
		a.logger.Errorf("error init sender: ", err)
		return err
	}
	return nil
}

func (a *Sender) Start(ctx context.Context) {
	a.Subscriber.Handle(ctx, a.numWorkers, a.reconnectTimeOut)
}

func (a *Sender) handler(ctx context.Context, contentEncoding string, content []byte) {
	log.Printf("%s %s", contentEncoding, content)
}

func NewAppSender(
	logg *logger.Logger,
	numWorkers int,
	dialString string,
	exchangeName string,
	routingKey string,
	queueName string,
	reconnectTimeOut time.Duration,
	subscriber queue.Subscriber) *Sender {
	return &Sender{
		logger:           logg,
		numWorkers:       numWorkers,
		dialString:       dialString,
		exchangeName:     exchangeName,
		routingKey:       routingKey,
		queueName:        queueName,
		reconnectTimeOut: reconnectTimeOut,
		Subscriber:       subscriber,
	}
}
