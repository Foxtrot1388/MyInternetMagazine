package config

import (
	"os"
	"sync"
)

type Config struct {
	Host       string
	Port       string
	User       string
	Pass       string
	SigningKey string
	KafkaHost  string
}

var (
	instance *Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		instance = &Config{
			Host:       getEnv("POSTGRES_HOST", "localhost"),
			Port:       getEnv("POSTGRES_PORT", "5432"),
			User:       getEnv("POSTGRES_USER", "root"),
			Pass:       getEnv("POSTGRES_PASSWORD", "root"),
			SigningKey: getEnv("SIGNING_KEY", ""),
			KafkaHost:  getEnv("KAFKA_HOST", ""),
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
