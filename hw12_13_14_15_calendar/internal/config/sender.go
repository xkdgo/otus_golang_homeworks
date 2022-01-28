package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type SenderConfig struct {
	Logger LoggerConf `mapstructure:"logger"`
	Sender SenderConf `mapstructure:"sender"`
	Broker BrokerConf `mapstructure:"broker"`
}

type SenderConf struct {
	RoutingKey       string `mapstructure:"routingkey"`
	ReconnectTimeOut string `mapstructure:"reconnectmsec"`
	NumWorkers       int    `mapstructure:"numworkers"`
}

func NewSenderConfig(cfgFile string, serviceName string) (SenderConfig, error) {
	viper.SetDefault("logger.level", defaultLoggerLevel)
	viper.SetDefault("broker.dialstring", defaultBrokerDialString)
	viper.SetDefault("sender.reconnectmsec", defaultReconnectTimeOut)
	viper.SetDefault("sender.routingkey", defaultRoutingKey)
	viper.SetDefault("sender.numworkers", defaultNumWorkers)
	var config SenderConfig
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return SenderConfig{}, fmt.Errorf("failed to read config: %w", err)
	}
	viper.SetEnvPrefix(serviceName)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	err := viper.Unmarshal(&config)
	if err != nil {
		return SenderConfig{}, fmt.Errorf("unable to decode into struct, %w", err)
	}
	fmt.Printf("%+v\n", config)
	fmt.Println(viper.AllSettings())
	return config, nil
}
