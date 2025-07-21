package ws

import "time"

// Option is a functional option for configuring a WebSocket handler.
type Option func(*WSConfig)

// WithRedis configures the handler to use a Redis message broker.
func WithRedis(config RedisConfig) Option {
	return func(c *WSConfig) {
		c.MessageBroker = MessageBrokerConfig{
			Type:  "redis",
			Redis: config,
		}
	}
}

// WithAutoSync enables or disables auto-sync functionality.
func WithAutoSync(enabled bool) Option {
	return func(c *WSConfig) {
		c.EnableAutoSync = enabled
	}
}

// WithSyncChannel sets the channel name for synchronization.
func WithSyncChannel(channel string) Option {
	return func(c *WSConfig) {
		c.SyncChannel = channel
	}
}

// WithPingInterval sets the interval for sending ping messages.
func WithPingInterval(interval time.Duration) Option {
	return func(c *WSConfig) {
		c.PingInterval = interval
	}
}

// WithPongWait sets the duration to wait for a pong response.
func WithPongWait(wait time.Duration) Option {
	return func(c *WSConfig) {
		c.PongWait = wait
	}
}

// WithWriteWait sets the duration to wait when writing a message.
func WithWriteWait(wait time.Duration) Option {
	return func(c *WSConfig) {
		c.WriteWait = wait
	}
}

// WithMaxMessageSize sets the maximum allowed message size.
func WithMaxMessageSize(size int64) Option {
	return func(c *WSConfig) {
		c.MaxMessageSize = size
	}
}

// WithBufferSize sets the buffer size for the send channel.
func WithBufferSize(size int) Option {
	return func(c *WSConfig) {
		c.BufferSize = size
	}
}
