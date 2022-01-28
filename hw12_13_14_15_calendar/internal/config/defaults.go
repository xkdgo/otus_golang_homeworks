package config

var (
	defaultLoggerLevel      = "INFO"
	defaultServerPort       = "8080"
	defaultStorageType      = "in-memory"
	defaultQuery            = "20s"
	defaultTTL              = "365"
	defaultBrokerDialString = "amqp://guest:guest@localhost:5672/"
	defaultReconnectTimeOut = "10000"
	defaultRoutingKey       = "calendar_sender"
	defaultNumWorkers       = 1
)
