package config

import (
	"fmt"
	"os"
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
	Logfile string `mapstructure:"logfile"`
	Level   string `mapstructure:"level"`
}

type ServerConf struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type StorageConf struct {
	Type string `mapstructure:"type"`
	DSN  string `mapstructure:"dsn"`
}

func NewCalendarConfigFromFile(cfgFile string, serviceName string) (CalendarConfig, error) {
	viper.SetDefault("server.port", defaultServerPort)
	viper.SetDefault("db.type", defaultStorageType)
	var config CalendarConfig
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return CalendarConfig{}, fmt.Errorf("failed to read config: %w", err)
	}
	// viper.SetEnvPrefix(serviceName)
	// replacer := strings.NewReplacer(".", "_")
	// viper.SetEnvKeyReplacer(replacer)
	// viper.AutomaticEnv()
	err := viper.Unmarshal(&config)
	if err != nil {
		return CalendarConfig{}, fmt.Errorf("unable to decode into struct, %w", err)
	}
	fmt.Printf("%+v\n", config)
	fmt.Println(viper.AllSettings())
	return config, nil
}

func NewCalendarConfigFromEnv(serviceName string) (CalendarConfig, error) {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("grpc.port", "9090")
	viper.SetDefault("db.type", "sql")
	user := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	viper.SetDefault("db.dsn", fmt.Sprintf("postgres://%s:%s@%s:%s/calendar?sslmode=disable", user, passwd, host, port))
	viper.SetDefault("logger.level", "INFO")
	var config CalendarConfig
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
