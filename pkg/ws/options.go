package ws

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Option is a functional option for configuring a WebSocket handler.
type Option func(*WSHandler)

// WithRedis is a convenience function that creates a Redis broker and sets it.
func WithRedis(client *redis.Client) Option {
	return func(h *WSHandler) {
		broker, err := NewRedisMessageBroker(client)
		if err != nil {
			// In a real-world app, you might want to return an error here
			// but for simplicity, we'll log and continue with the no-op broker.
			log.Printf("Could not create Redis message broker: %v", err)
			return
		}
		h.broker = broker
	}
}

// WithMessageBroker sets the message broker for the WebSocket handler.
func WithMessageBroker(broker MessageBroker) Option {
	return func(h *WSHandler) {
		if broker != nil {
			h.broker = broker
		}
	}
}

// WithAutoSync enables or disables auto-sync functionality.
func WithAutoSync(enabled bool) Option {
	return func(h *WSHandler) {
		h.config.EnableAutoSync = enabled
	}
}

// WithSyncChannel sets the channel name for synchronization.
func WithSyncChannel(channel string) Option {
	return func(h *WSHandler) {
		h.config.SyncChannel = channel
	}
}

// WithPingInterval sets the interval for sending ping messages.
func WithPingInterval(interval time.Duration) Option {
	return func(h *WSHandler) {
		h.config.PingInterval = interval
	}
}

// WithPongWait sets the duration to wait for a pong response.
func WithPongWait(wait time.Duration) Option {
	return func(h *WSHandler) {
		h.config.PongWait = wait
	}
}

// WithWriteWait sets the duration to wait when writing a message.
func WithWriteWait(wait time.Duration) Option {
	return func(h *WSHandler) {
		h.config.WriteWait = wait
	}
}

// WithMaxMessageSize sets the maximum allowed message size.
func WithMaxMessageSize(size int64) Option {
	return func(h *WSHandler) {
		h.config.MaxMessageSize = size
	}
}

// WithBufferSize sets the buffer size for the send channel.
func WithBufferSize(size int) Option {
	return func(h *WSHandler) {
		h.config.BufferSize = size
	}
}
