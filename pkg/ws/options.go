package ws

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Option is a functional option for configuring a WebSocket handler.
type Option func(*ConnectionManager)

// WithRedis is a convenience function that creates a Redis broker and sets it.
func WithRedis(client *redis.Client) Option {
	return func(cm *ConnectionManager) {
		broker, err := NewRedisMessageBroker(client)
		if err != nil {
			// In a real-world app, you might want to return an error here
			// but for simplicity, we'll log and continue with the no-op broker.
			log.Printf("Could not create Redis message broker: %v", err)
			return
		}
		cm.broker = broker
	}
}

// WithMessageBroker sets the message broker for the WebSocket handler.
func WithMessageBroker(broker MessageBroker) Option {
	return func(cm *ConnectionManager) {
		if broker != nil {
			cm.broker = broker
		}
	}
}

// WithAutoSync enables or disables auto-sync functionality.
func WithAutoSync(enabled bool) Option {
	return func(cm *ConnectionManager) {
		cm.config.EnableAutoSync = enabled
	}
}

// WithSyncChannel sets the channel name for synchronization.
func WithSyncChannel(channel string) Option {
	return func(cm *ConnectionManager) {
		cm.config.SyncChannel = channel
	}
}

// WithPingInterval sets the interval for sending ping messages.
func WithPingInterval(interval time.Duration) Option {
	return func(cm *ConnectionManager) {
		cm.config.PingInterval = interval
	}
}

// WithPongWait sets the duration to wait for a pong response.
func WithPongWait(wait time.Duration) Option {
	return func(cm *ConnectionManager) {
		cm.config.PongWait = wait
	}
}

// WithWriteWait sets the duration to wait when writing a message.
func WithWriteWait(wait time.Duration) Option {
	return func(cm *ConnectionManager) {
		cm.config.WriteWait = wait
	}
}

// WithMaxMessageSize sets the maximum allowed message size.
func WithMaxMessageSize(size int64) Option {
	return func(cm *ConnectionManager) {
		cm.config.MaxMessageSize = size
	}
}

// WithBufferSize sets the buffer size for the send channel.
func WithBufferSize(size int) Option {
	return func(cm *ConnectionManager) {
		cm.config.BufferSize = size
	}
}
