package queue

type Notifier interface {
	Init() error
	Publish(routingKey string, contentType string, body []byte) error
	Stop() error
	Listen()
}
