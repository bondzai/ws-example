package ws

// MessageBrokerConfig holds configuration for different message broker types
type MessageBrokerConfig struct {
	Type  string `mapstructure:"type"` // "redis", "kafka", "nats", etc.
	Redis RedisConfig
	Kafka KafkaConfig
	Nats  NatsConfig
}

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	URL      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// KafkaConfig holds Kafka-specific configuration
type KafkaConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	Topic    string   `mapstructure:"topic"`
	GroupID  string   `mapstructure:"group_id"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

// NatsConfig holds NATS-specific configuration
type NatsConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Cluster  string `mapstructure:"cluster"`
}
