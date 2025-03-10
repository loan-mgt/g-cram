package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	ServerPort int
	LogLevel   string
}

// New returns a new Config with values from environment variables or defaults
func New() *Config {
	return &Config{
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as int with a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
