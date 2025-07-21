package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application's configuration.
type Config struct {
	HttpPort string      `mapstructure:"http_port"`
	BaseUrl  string      `mapstructure:"base_url"`
	Redis    RedisConfig `mapstructure:"redis"`
}

// RedisConfig holds Redis-specific connection details.
type RedisConfig struct {
	URL      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// NewConfig initializes and loads the application configuration.
func NewConfig() *Config {
	// --- Set Default Values ---
	viper.SetDefault("http_port", "8080")
	viper.SetDefault("base_url", "http://localhost:8080")
	viper.SetDefault("redis.url", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// --- Load Configuration from .env and Environment Variables ---
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the config file if it exists.
	_ = viper.ReadInConfig()

	// Unmarshal the configuration into the Config struct.
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("unable to decode config into struct: %v", err))
	}

	return &cfg
}
