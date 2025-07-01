package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbHost    string
	DbUser    string
	DbPass    string
	DbName    string
	DbPort    string
	DbSslMode string
	DbDsn     string

	RedisUrl string

	HttpPort string
	BaseUrl  string
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	return &Config{
		DbHost:    viper.GetString("DB_HOST"),
		DbUser:    viper.GetString("DB_USER"),
		DbPass:    viper.GetString("DB_PASS"),
		DbName:    viper.GetString("DB_NAME"),
		DbPort:    viper.GetString("DB_PORT"),
		DbSslMode: viper.GetString("DB_SSLMODE"),
		DbDsn:     viper.GetString("DB_DSN"),
		RedisUrl:  viper.GetString("REDIS_URL"),
		HttpPort:  viper.GetString("HTTP_PORT"),
		BaseUrl:   viper.GetString("BASE_URL"),
	}
}
