package rabbit

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/queue"
)

var _ queue.Notifier = (*Sender)(nil)

type Sender struct {
	exchange   Exchange
	dialString string
	log        *logger.Logger
	connection *amqp.Connection
	channel    *amqp.Channel
	errCh      chan *amqp.Error
}

type Exchange struct {
	Name    string
	Type    string
	Durable bool
}

func NewExchange(exchange, typ string, durable bool) Exchange {
	return Exchange{
		Name:    exchange,
		Type:    typ,
		Durable: durable,
	}
}

func (s *Sender) Init() error {
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
	err := s.ChannelClose()
	if err != nil {
		s.log.Errorf("unable close channel", err)
	}
	if s.connection == nil {
		return errors.New("connection is nil")
	}
	s.connection.Close()
	return nil
}

func (s *Sender) Publish(routingKey string, contentType string, body []byte) error {
	message := amqp.Publishing{
		ContentType: contentType,
		Body:        body,
	}
	if s.channel == nil {
		return errors.New("channel is nil")
	}
	return s.channel.Publish(s.exchange.Name, routingKey, false, false, message)
}

func NewSender(exchangeName, typ string, durable bool, dialString string, logg *logger.Logger) (*Sender, error) {
	s := &Sender{
		exchange:   NewExchange(exchangeName, typ, durable),
		dialString: dialString,
		log:        logg,
		connection: nil,
		channel:    nil,
	}
	err := s.Init()
	if err != nil {
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
		err := s.Init()
		for err != nil {
			err = s.Init()
		}
		s.log.Infof("reconnected to server %s", s.dialString)
	}
}
