package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds the application's configuration.
type Config struct {
	AppName  string
	HttpPort string
	BaseUrl  string
	Redis    RedisConfig
	Mongo    MongoConfig
}

// RedisConfig holds Redis-specific connection details.
type RedisConfig struct {
	URI      string
	Password string
	DB       int
}

// MongoConfig holds MongoDB-specific connection details.
type MongoConfig struct {
	URI      string
	Database string
}

// NewConfig initializes and loads the application configuration by explicitly
// reading each key from the environment.
func NewConfig() *Config {
	// --- Load Configuration from .env and Environment Variables ---
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Read the config file if it exists.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Info: .env file not found. Using environment variables.")
		} else {
			log.Printf("Warning: Could not read .env file: %s", err)
		}
	}

	// This allows Viper to read environment variables directly.
	viper.AutomaticEnv()

	// --- Explicitly Read and Populate Config ---
	cfg := &Config{
		AppName:  getEnv("APP_NAME", "WebSocketApp"),
		HttpPort: getEnv("HTTP_PORT", "8080"),
		BaseUrl:  getEnv("BASE_URL", "http://localhost:8080"),
		Redis: RedisConfig{
			URI:      getEnv("REDIS_URL", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       viper.GetInt("REDIS_DB"), // Defaults to 0 if not set
		},
		Mongo: MongoConfig{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGO_DATABASE", "chat_db"),
		},
	}

	return cfg
}

// getEnv is a helper to read an environment variable or return a default value.
func getEnv(key, defaultValue string) string {
	if value := viper.GetString(key); value != "" {
		return value
	}
	return defaultValue
}
