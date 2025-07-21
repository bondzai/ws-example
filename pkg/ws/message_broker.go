package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// MessageBroker defines the interface for different message broker implementations
type MessageBroker interface {
	// Publish sends a message to the broker
	Publish(ctx context.Context, topic string, message interface{}) error
	// Subscribe starts listening for messages on a topic
	Subscribe(ctx context.Context, topic string, handler func(message []byte)) error
	// Close closes the broker connection
	Close() error
	// GetType returns the broker type
	GetType() string
}

// BaseMessageBroker provides common functionality for all message brokers
type BaseMessageBroker struct {
	ctx      context.Context
	cancel   context.CancelFunc
	handlers map[string]func(message []byte)
}

// NewBaseMessageBroker creates a new base message broker
func NewBaseMessageBroker() *BaseMessageBroker {
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseMessageBroker{
		ctx:      ctx,
		cancel:   cancel,
		handlers: make(map[string]func(message []byte)),
	}
}

// RegisterHandler registers a message handler for a topic
func (b *BaseMessageBroker) RegisterHandler(topic string, handler func(message []byte)) {
	b.handlers[topic] = handler
}

// GetHandler returns the handler for a topic
func (b *BaseMessageBroker) GetHandler(topic string) (func(message []byte), bool) {
	handler, exists := b.handlers[topic]
	return handler, exists
}

// GetContext returns the broker context
func (b *BaseMessageBroker) GetContext() context.Context {
	return b.ctx
}

// Cancel cancels the broker context
func (b *BaseMessageBroker) Cancel() {
	b.cancel()
}

// RedisMessageBroker implements MessageBroker for Redis Pub/Sub
type RedisMessageBroker struct {
	*BaseMessageBroker
	client *redis.Client
}

// NewRedisMessageBroker creates a new Redis message broker
func NewRedisMessageBroker(config RedisConfig) (*RedisMessageBroker, error) {
	base := NewBaseMessageBroker()

	client := redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	// Test the connection
	if err := client.Ping(base.GetContext()).Err(); err != nil {
		base.Cancel()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
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
	defer pubsub.Close()

	// Start listening for messages
	go func() {
		for {
			select {
			case <-r.GetContext().Done():
				return
			default:
				msg, err := pubsub.ReceiveMessage(ctx)
				if err != nil {
					log.Printf("Redis subscription error: %v", err)
					continue
				}

				if handler, exists := r.GetHandler(topic); exists {
					handler([]byte(msg.Payload))
				}
			}
		}
	}()

	return nil
}

func (r *RedisMessageBroker) Close() error {
	r.Cancel()
	return r.client.Close()
}

func (r *RedisMessageBroker) GetType() string {
	return "redis"
}

// NoOpMessageBroker is a no-operation message broker for testing or when no broker is needed
type NoOpMessageBroker struct {
	*BaseMessageBroker
}

func NewNoOpMessageBroker() *NoOpMessageBroker {
	return &NoOpMessageBroker{
		BaseMessageBroker: NewBaseMessageBroker(),
	}
}

func (n *NoOpMessageBroker) Publish(ctx context.Context, topic string, message interface{}) error {
	return nil
}

func (n *NoOpMessageBroker) Subscribe(ctx context.Context, topic string, handler func(message []byte)) error {
	return nil
}

func (n *NoOpMessageBroker) Close() error {
	n.Cancel()
	return nil
}

func (n *NoOpMessageBroker) GetType() string {
	return "noop"
}

// MessageBrokerFactory creates message brokers based on configuration
type MessageBrokerFactory struct{}

// NewMessageBrokerFactory creates a new message broker factory
func NewMessageBrokerFactory() *MessageBrokerFactory {
	return &MessageBrokerFactory{}
}

// CreateMessageBroker creates a message broker based on the provided configuration
func (f *MessageBrokerFactory) CreateMessageBroker(config MessageBrokerConfig) (MessageBroker, error) {
	switch config.Type {
	case "redis":
		return NewRedisMessageBroker(config.Redis)
	case "noop", "":
		return NewNoOpMessageBroker(), nil
	default:
		return nil, fmt.Errorf("unsupported message broker type: %s", config.Type)
	}
}

// MessageBrokerManager manages the lifecycle of message brokers
type MessageBrokerManager struct {
	broker MessageBroker
	config MessageBrokerConfig
}

// NewMessageBrokerManager creates a new message broker manager
func NewMessageBrokerManager(config MessageBrokerConfig) (*MessageBrokerManager, error) {
	factory := NewMessageBrokerFactory()
	broker, err := factory.CreateMessageBroker(config)
	if err != nil {
		return nil, err
	}

	return &MessageBrokerManager{
		broker: broker,
		config: config,
	}, nil
}

// GetBroker returns the managed message broker
func (m *MessageBrokerManager) GetBroker() MessageBroker {
	return m.broker
}

// Close closes the message broker
func (m *MessageBrokerManager) Close() error {
	return m.broker.Close()
}

// Publish publishes a message to the broker
func (m *MessageBrokerManager) Publish(ctx context.Context, topic string, message interface{}) error {
	return m.broker.Publish(ctx, topic, message)
}

// Subscribe subscribes to messages from the broker
func (m *MessageBrokerManager) Subscribe(ctx context.Context, topic string, handler func(message []byte)) error {
	return m.broker.Subscribe(ctx, topic, handler)
}
