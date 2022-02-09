package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type SchedulerConfig struct {
	Logger    LoggerConf    `mapstructure:"logger"`
	Storage   StorageConf   `mapstructure:"db"`
	Scheduler SchedulerTune `mapstructure:"scheduler"`
	Broker    BrokerConf    `mapstructure:"broker"`
}

type SchedulerTune struct {
	TimeoutQuery     string `mapstructure:"timeoutquery"`
	TTL              string `mapstructure:"ttldays"`
	ReconnectTimeOut string `mapstructure:"reconnectmsec"`
}

type BrokerConf struct {
	DialString string `mapstructure:"dialstring"`
}

func NewSchedulerConfigFromFile(cfgFile string, serviceName string) (SchedulerConfig, error) {
	viper.SetDefault("db.type", defaultStorageType)
	viper.SetDefault("scheduler.timeoutquery", defaultQuery)
	viper.SetDefault("scheduler.ttldays", defaultTTL)
	viper.SetDefault("broker.dialstring", defaultBrokerDialString)
	viper.SetDefault("scheduler.reconnectmsec", defaultReconnectTimeOut)
	var config SchedulerConfig
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return SchedulerConfig{}, fmt.Errorf("failed to read config: %w", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		return SchedulerConfig{}, fmt.Errorf("unable to decode into struct, %w", err)
	}
	fmt.Printf("%+v\n", config)
	fmt.Println(viper.AllSettings())
	return config, nil
}

func NewSchedulerConfigFromEnv(serviceName string) (SchedulerConfig, error) {
	viper.SetDefault("scheduler.timeoutquery", defaultQuery)
	viper.SetDefault("scheduler.ttldays", defaultTTL)
	viper.SetDefault("scheduler.reconnectmsec", defaultReconnectTimeOut)
	viper.SetDefault("db.type", "sql")
	viper.SetDefault("db.dsn", buildDsnFromEnv())
	viper.SetDefault("logger.level", defaultLoggerLevel)
	viper.SetDefault("broker.dialstring", buildRmqDialFromEnv())
	var config SchedulerConfig
	viper.SetEnvPrefix(serviceName)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	err := viper.Unmarshal(&config)
	if err != nil {
		return SchedulerConfig{}, fmt.Errorf("unable to decode into struct, %w", err)
	}
	fmt.Printf("%+v\n", config)
	fmt.Println(viper.AllSettings())
	return config, nil
}
