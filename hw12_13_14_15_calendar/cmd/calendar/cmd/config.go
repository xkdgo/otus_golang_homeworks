package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

// var ErrFailToReadConfig = errors.New("failed to read config")

type Config struct {
	Logger LoggerConf `mapstructure:"logger"`
	// TODO
}

type LoggerConf struct {
	Level     string `mapstructure:"level"`
	SomeParam string `mapstructure:"someparam"`
	// TODO
}

func NewConfig(cfgFile string) (Config, error) {
	viper.SetDefault("logger.someparam", "defaultString")
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
	test := viper.GetString("logger.someparam")
	fmt.Println(test)
	fmt.Println(viper.AllSettings())
	return config, nil
}
