package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type CalendarConfig struct {
	Logger     LoggerConf  `mapstructure:"logger"`
	ServerHTTP ServerConf  `mapstructure:"server"`
	ServerGRPC ServerConf  `mapstructure:"grpc"`
	Storage    StorageConf `mapstructure:"db"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
}

type ServerConf struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type StorageConf struct {
	Type string `mapstructure:"type"`
	DSN  string `mapstructure:"dsn"`
}

func NewCalendarConfig(cfgFile string, serviceName string) (CalendarConfig, error) {
	viper.SetDefault("server.port", defaultServerPort)
	viper.SetDefault("db.type", defaultStorageType)
	var config CalendarConfig
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return CalendarConfig{}, fmt.Errorf("failed to read config: %w", err)
	}
	viper.SetEnvPrefix(serviceName)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	err := viper.Unmarshal(&config)
	if err != nil {
		return CalendarConfig{}, fmt.Errorf("unable to decode into struct, %w", err)
	}
	fmt.Printf("%+v\n", config)
	fmt.Println(viper.AllSettings())
	return config, nil
}
