package queue

type Notifier interface {
	Init(opts ...Option) error
	Publish(routingKey, contentType string, body []byte) error
	Stop() error
	Listen()
}
