package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type SchedulerConfig struct {
	Logger    LoggerConf    `mapstructure:"logger"`
	Storage   StorageConf   `mapstructure:"db"`
	Scheduler SchedulerConf `mapstructure:"scheduler"`
}

type SchedulerConf struct {
	TimeoutQuery string `mapstructure:"timeoutquery"`
}

func NewSchedulerConfig(cfgFile string, serviceName string) (SchedulerConfig, error) {
	viper.SetDefault("db.type", defaultStorageType)
	viper.SetDefault("scheduler.timeoutquery", defaultQuery)
	var config SchedulerConfig
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return SchedulerConfig{}, fmt.Errorf("failed to read config: %w", err)
	}
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
