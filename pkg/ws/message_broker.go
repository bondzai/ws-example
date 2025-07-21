package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// MessageBroker defines the interface for different message broker implementations.
type MessageBroker interface {
	Publish(ctx context.Context, topic string, message interface{}) error
	Subscribe(ctx context.Context, topic string, handler func(message []byte)) error
	Close() error
	GetType() string
}

// BaseMessageBroker provides common functionality for all message brokers.
type BaseMessageBroker struct {
	ctx      context.Context
	cancel   context.CancelFunc
	handlers map[string]func(message []byte)
}

// NewBaseMessageBroker creates a new base message broker.
func NewBaseMessageBroker() *BaseMessageBroker {
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseMessageBroker{
		ctx:      ctx,
		cancel:   cancel,
		handlers: make(map[string]func(message []byte)),
	}
}

// RegisterHandler registers a message handler for a topic.
func (b *BaseMessageBroker) RegisterHandler(topic string, handler func(message []byte)) {
	b.handlers[topic] = handler
}

// GetHandler returns the handler for a topic.
func (b *BaseMessageBroker) GetHandler(topic string) (func(message []byte), bool) {
	handler, exists := b.handlers[topic]
	return handler, exists
}

// GetContext returns the broker context.
func (b *BaseMessageBroker) GetContext() context.Context {
	return b.ctx
}

// Cancel cancels the broker context.
func (b *BaseMessageBroker) Cancel() {
	b.cancel()
}

// RedisMessageBroker implements MessageBroker for Redis Pub/Sub.
type RedisMessageBroker struct {
	*BaseMessageBroker
	client *redis.Client
}

// NewRedisMessageBroker creates a new Redis message broker from an existing client.
// It allows injecting a pre-configured redis.Client.
func NewRedisMessageBroker(client *redis.Client) (MessageBroker, error) {
	base := NewBaseMessageBroker()

	// Test the connection
	if err := client.Ping(base.GetContext()).Err(); err != nil {
		base.Cancel()
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return &RedisMessageBroker{
		BaseMessageBroker: base,
		client:            client,
	}, nil
}

func (r *RedisMessageBroker) Publish(ctx context.Context, topic string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	return r.client.Publish(ctx, topic, data).Err()
}

func (r *RedisMessageBroker) Subscribe(ctx context.Context, topic string, handler func(message []byte)) error {
	r.RegisterHandler(topic, handler)

	pubsub := r.client.Subscribe(ctx, topic)

	// Start listening for messages in a separate goroutine.
	go func() {
		defer pubsub.Close()
		for {
			select {
			case <-r.GetContext().Done():
				// Context was cancelled, stop the goroutine.
				return
			case msg, ok := <-pubsub.Channel():
				if !ok {
					// Channel was closed.
					return
				}
				if handler, exists := r.GetHandler(topic); exists {
					handler([]byte(msg.Payload))
				}
			}
		}
	}()

	return nil
}

// Close cancels the context for the message broker's goroutines.
// It does not close the Redis client, as its lifecycle is managed externally.
func (r *RedisMessageBroker) Close() error {
	r.Cancel()
	return nil
}

func (r *RedisMessageBroker) GetType() string {
	return "redis"
}

// NoOpMessageBroker is a no-operation message broker for testing or when no broker is needed.
type NoOpMessageBroker struct {
	*BaseMessageBroker
}

// NewNoOpMessageBroker creates a new no-op message broker.
func NewNoOpMessageBroker() MessageBroker {
	return &NoOpMessageBroker{
		BaseMessageBroker: NewBaseMessageBroker(),
	}
}

func (n *NoOpMessageBroker) Publish(ctx context.Context, topic string, message interface{}) error {
	// This broker does nothing.
	return nil
}

func (n *NoOpMessageBroker) Subscribe(ctx context.Context, topic string, handler func(message []byte)) error {
	// This broker does nothing.
	return nil
}

// Close cancels the context for the no-op broker.
func (n *NoOpMessageBroker) Close() error {
	n.Cancel()
	return nil
}

func (n *NoOpMessageBroker) GetType() string {
	return "noop"
}
