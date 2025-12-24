package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Port     string
	LogLevel string
	MaxLimit int
}

// Load reads configuration from environment
func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		MaxLimit: getEnvAsInt("MAX_LIMIT", 10000),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
