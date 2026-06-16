package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort string

	DBHost string
	DBPort string

	DBUser     string
	DBPassword string
	DBName     string
}

func Load() *Config {

	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),

		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),

		DBUser:     getEnv("DB_USER", "wallet"),
		DBPassword: getEnv("DB_PASSWORD", "wallet"),
		DBName:     getEnv("DB_NAME", "wallet"),
	}

	validate(cfg)

	return cfg
}

func validate(
	cfg *Config,
) {

	if cfg.DBHost == "" {
		log.Fatal("DB_HOST missing")
	}
}

func getEnv(
	key string,
	defaultValue string,
) string {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
