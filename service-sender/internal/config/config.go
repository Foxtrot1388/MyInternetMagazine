package config

import (
	"os"
	"sync"
)

type Config struct {
	Host       string
	EmailTopic string
	Email      EmailConfig
}

type EmailConfig struct {
	Sender            string
	Pass              string
	SmtpAuthAddress   string
	SmtpServerAddress string
}

var (
	instance *Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		instance = &Config{
			Host:       getEnv("KAFKA_HOST", "localhost"),
			EmailTopic: getEnv("KAFKA_TOPIC_EMAIL", "email"),
			Email: EmailConfig{
				Sender:            getEnv("SENDER_EMAIL", ""),
				Pass:              getEnv("PASS_EMAIL", ""),
				SmtpAuthAddress:   getEnv("SMTP_AUTH_ADDRESS_EMAIL", ""),
				SmtpServerAddress: getEnv("SMTP_SERVER_EMAIL", ""),
			},
		}
	})
	return instance
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
