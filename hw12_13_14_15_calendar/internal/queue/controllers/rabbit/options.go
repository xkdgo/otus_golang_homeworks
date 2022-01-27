package rabbit

import "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/queue"

type exchangeNameKey struct{}

func WithExchangeName(exchangeName string) queue.Option {
	return queue.SetOption(exchangeNameKey{}, exchangeName)
}

type exchangeTypeKey struct{}

func WithExchangeType(exchangeType string) queue.Option {
	return queue.SetOption(exchangeTypeKey{}, exchangeType)
}

type exchangeDurableKey struct{}

func WithExchangeDurable(exchangeDurable bool) queue.Option {
	return queue.SetOption(exchangeDurableKey{}, exchangeDurable)
}

// type routingKeyKey struct{}

// func WithRoutingKey(routingKey string) queue.Option {
// 	return queue.SetOption(routingKeyKey{}, routingKey)
// }
