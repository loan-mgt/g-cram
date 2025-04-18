package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	ServerPort      int
	LogLevel        string
	RabbitMQHost    string
	RabbitMQPort    int
	RabbitMQUser    string
	RabbitMQPass    string
	DBPath          string
	FrontUrl        string
	FrontDomain     string
	ClientSecret    string
	ClientID        string
	TokenURI        string
	VAPIDPublicKey  string
	VAPIDPrivateKey string
	VAPIDSubject    string
	Salt            string
}

// New returns a new Config with values from environment variables or defaults
func New() *Config {
	return &Config{
		ServerPort:      getEnvAsInt("SERVER_PORT", 8080),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		RabbitMQHost:    getEnv("RABBITMQ_HOST", "localhost"),
		RabbitMQPort:    getEnvAsInt("RABBITMQ_PORT", 5672),
		RabbitMQUser:    getEnv("RABBITMQ_USER", "guest"),
		RabbitMQPass:    getEnv("RABBITMQ_PASS", "guest"),
		DBPath:          getEnv("DB_PATH", "/database/g-cram.db"),
		FrontUrl:        getEnv("FRONT_URL", "http://localhost:8080"),
		FrontDomain:     getEnv("FRONT_DOMAIN", "localhost"),
		ClientSecret:    getEnv("CLIENT_SECRET", ""),
		ClientID:        getEnv("CLIENT_ID", ""),
		TokenURI:        getEnv("TOKEN_URI", "https://oauth2.googleapis.com/token"),
		VAPIDPublicKey:  getEnv("VAPID_PUBLIC_KEY", ""),
		VAPIDPrivateKey: getEnv("VAPID_PRIVATE_KEY", ""),
		VAPIDSubject:    getEnv("VAPID_SUBJECT", ""),
		Salt:            getEnv("SALT", ""),
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
