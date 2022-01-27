package rabbit

import (
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/queue"
)

var _ queue.Notifier = (*Sender)(nil)

type Publisher struct {
	exchangeName string
	routingKey   string
	channel      *amqp.Channel
}

func (p *Publisher) Publish(exchangeName, routingKey string, message amqp.Publishing) error {
	return p.channel.Publish(p.exchangeName, p.routingKey, false, false, message)
}

func NewPublisher(exchangeName, routingKey string, ch *amqp.Channel) *Publisher {
	return &Publisher{
		exchangeName: exchangeName,
		routingKey:   routingKey,
		channel:      ch,
	}
}

type Exchange struct {
	Name    string
	Type    string
	Durable bool
}

var DefaultExchange = Exchange{
	Name:    "exchange",
	Type:    "direct",
	Durable: true,
}

type Sender struct {
	opts             queue.Options
	exchange         Exchange
	publishers       map[string]*Publisher
	dialString       string
	reconnectTimeOut time.Duration
	log              *logger.Logger
	connection       *amqp.Connection
	channel          *amqp.Channel
	errCh            chan *amqp.Error
}

func (s *Sender) Init(opts ...queue.Option) error {
	for _, o := range opts {
		o(&s.opts)
	}
	if exchangeName, ok := s.opts.Context.Value(exchangeNameKey{}).(string); ok {
		s.exchange.Name = exchangeName
	}
	if exchangeType, ok := s.opts.Context.Value(exchangeTypeKey{}).(string); ok {
		s.exchange.Type = exchangeType
	}
	if exchangeDurable, ok := s.opts.Context.Value(exchangeDurableKey{}).(bool); ok {
		s.exchange.Durable = exchangeDurable
	}
	err := s.Dial()
	if err != nil {
		return err
	}
	err = s.ConnectExchange()
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) reconnect() error {
	err := s.Dial()
	if err != nil {
		return err
	}
	err = s.ConnectExchange()
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) Dial() error {
	conn, err := amqp.Dial(s.dialString)
	if err != nil {
		return err
	}
	s.connection = conn
	return nil
}

func (s *Sender) ConnectExchange() error {
	var err error
	if s.channel == nil || s.channel.IsClosed() {
		s.channel, err = s.connection.Channel()
	}
	if err != nil {
		return err
	}
	err = s.channel.ExchangeDeclare(
		s.exchange.Name,    // name
		s.exchange.Type,    // type
		s.exchange.Durable, // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) ChannelClose() error {
	if s.channel == nil {
		return errors.New("channel is nil")
	}
	return s.channel.Close()
}

func (s *Sender) Stop() error {
	if err := s.ChannelClose(); err != nil {
		s.log.Errorf("unable close channel", err)
	}
	if s.connection == nil {
		return errors.New("connection is nil")
	}
	s.connection.Close()
	return nil
}

func (s *Sender) Publish(routingKey, contentType string, body []byte) error {
	message := amqp.Publishing{
		ContentType: contentType,
		Body:        body,
	}
	if s.channel == nil {
		return errors.New("channel is nil")
	}
	_, ok := s.publishers[routingKey]
	if !ok {
		s.publishers[routingKey] = NewPublisher(s.exchange.Name, routingKey, s.channel)
	}
	return s.publishers[routingKey].Publish(s.exchange.Name, routingKey, message)
}

func NewSender(
	exchangeName, typ string,
	durable bool,
	dialString string,
	reconnectTimeOut time.Duration,
	logg *logger.Logger) (*Sender, error) {
	s := &Sender{
		exchange:         DefaultExchange,
		publishers:       make(map[string]*Publisher),
		dialString:       dialString,
		reconnectTimeOut: reconnectTimeOut,
		log:              logg,
		connection:       nil,
		channel:          nil,
	}
	if err := s.Init(
		WithExchangeName(exchangeName),
		WithExchangeType(typ),
		WithExchangeDurable(durable),
	); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Sender) Listen() {
	for {
		s.errCh = make(chan *amqp.Error)
		s.errCh = s.channel.NotifyClose(s.errCh)
		amqperr := <-s.errCh
		if amqperr == nil {
			s.log.Infof("listen catch nil error %v and exit", amqperr)
			return
		}
		s.log.Infof("listen catch error %v", amqperr)
		s.publishers = make(map[string]*Publisher)
		err := s.reconnect()
		for err != nil {
			time.Sleep(s.reconnectTimeOut)
			err = s.reconnect()
		}
		s.log.Infof("reconnected to server %s", s.dialString)
	}
}
