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

	// WebSocket Configuration
	WebSocket WebSocketConfig

	// Message Broker Configuration
	MessageBroker MessageBrokerConfig
}

type WebSocketConfig struct {
	PingInterval   int    `mapstructure:"ping_interval"`    // in seconds
	PongWait       int    `mapstructure:"pong_wait"`        // in seconds
	WriteWait      int    `mapstructure:"write_wait"`       // in seconds
	MaxMessageSize int64  `mapstructure:"max_message_size"` // in bytes
	BufferSize     int    `mapstructure:"buffer_size"`
	EnableAutoSync bool   `mapstructure:"enable_auto_sync"`
	SyncChannel    string `mapstructure:"sync_channel"`
}

type MessageBrokerConfig struct {
	Type  string `mapstructure:"type"` // "redis", "kafka", "nats", etc.
	Redis RedisConfig
	Kafka KafkaConfig
	Nats  NatsConfig
}

type RedisConfig struct {
	URL      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type KafkaConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	Topic    string   `mapstructure:"topic"`
	GroupID  string   `mapstructure:"group_id"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

type NatsConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Cluster  string `mapstructure:"cluster"`
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
		WebSocket: WebSocketConfig{
			PingInterval:   viper.GetInt("WS_PING_INTERVAL"),
			PongWait:       viper.GetInt("WS_PONG_WAIT"),
			WriteWait:      viper.GetInt("WS_WRITE_WAIT"),
			MaxMessageSize: viper.GetInt64("WS_MAX_MESSAGE_SIZE"),
			BufferSize:     viper.GetInt("WS_BUFFER_SIZE"),
			EnableAutoSync: viper.GetBool("WS_ENABLE_AUTO_SYNC"),
			SyncChannel:    viper.GetString("WS_SYNC_CHANNEL"),
		},
		MessageBroker: MessageBrokerConfig{
			Type: viper.GetString("MESSAGE_BROKER_TYPE"),
			Redis: RedisConfig{
				URL:      viper.GetString("REDIS_URL"),
				Password: viper.GetString("REDIS_PASSWORD"),
				DB:       viper.GetInt("REDIS_DB"),
				PoolSize: viper.GetInt("REDIS_POOL_SIZE"),
			},
			Kafka: KafkaConfig{
				Brokers:  viper.GetStringSlice("KAFKA_BROKERS"),
				Topic:    viper.GetString("KAFKA_TOPIC"),
				GroupID:  viper.GetString("KAFKA_GROUP_ID"),
				Username: viper.GetString("KAFKA_USERNAME"),
				Password: viper.GetString("KAFKA_PASSWORD"),
			},
			Nats: NatsConfig{
				URL:      viper.GetString("NATS_URL"),
				Username: viper.GetString("NATS_USERNAME"),
				Password: viper.GetString("NATS_PASSWORD"),
				Cluster:  viper.GetString("NATS_CLUSTER"),
			},
		},
	}
}
