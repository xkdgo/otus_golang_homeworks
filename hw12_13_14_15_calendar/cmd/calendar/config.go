package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf `mapstructure:"logger"`
	// TODO
}

type LoggerConf struct {
	Level     string `mapstructure:"level"`
	SomeParam string `mapstructure:"someparam"`
	// TODO
}

func NewConfig() Config {
	viper.SetDefault("logger.someparam", "defaultString")
	var C Config
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	err := viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	fmt.Printf("%+v\n", C)
	test := viper.GetString("logger.someparam")
	fmt.Println(test)
	fmt.Println(viper.AllSettings())
	return C
}

// TODO
