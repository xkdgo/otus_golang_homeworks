package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

var defaultServerPort = "8080"

type Config struct {
	Logger LoggerConf `mapstructure:"logger"`
	Server ServerConf `mapstructure:"server"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
}

type ServerConf struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func NewConfig(cfgFile string) (Config, error) {
	viper.SetDefault("server.port", defaultServerPort)
	var config Config
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("unable to decode into struct, %w", err)
	}
	fmt.Printf("%+v\n", config)
	fmt.Println(viper.AllSettings())
	return config, nil
}
