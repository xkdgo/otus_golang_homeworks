package config

import (
	"fmt"
	"os"
)

func buildDsnFromEnv() string {
	user := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/calendar?sslmode=disable", user, passwd, host, port)
}

func buildRmqDialFromEnv() string {
	rabbitUser := os.Getenv("RABBITMQ_DEFAULT_USER")
	rabbitPasswd := os.Getenv("RABBITMQ_DEFAULT_PASS")
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")
	if rabbitUser == "" {
		rabbitUser = "guest"
	}
	if rabbitPasswd == "" {
		rabbitPasswd = "guest"
	}
	if rabbitHost == "" {
		rabbitHost = "localhost"
	}
	if rabbitPort == "" {
		rabbitPort = "5672"
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%s",
		rabbitUser, rabbitPasswd, rabbitHost, rabbitPort)
}
