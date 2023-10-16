package config

import (
	"os"
	"sync"
)

type Config struct {
	Host  string
	Port  string
	User  string
	Pass  string
	Cashe CasheConfig
}

type CasheConfig struct {
	CasheHost string
	CashePort string
}

var (
	instance *Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		instance = &Config{
			Host: getEnv("POSTGRES_HOST", "localhost"),
			Port: getEnv("POSTGRES_PORT", "5433"),
			User: getEnv("POSTGRES_USER", "root"),
			Pass: getEnv("POSTGRES_PASSWORD", "root"),
			Cashe: CasheConfig{
				CasheHost: getEnv("REDIS_HOST", "localhost"),
				CashePort: getEnv("REDIS_PORT", "6379"),
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
