// pkg/kafka/types.go
package kafka

import (
	"context"
	"time"
)

// Producer is the interface for sending messages to Kafka.
type Producer interface {
	// Produce sends a message to the specified topic with given key and value.
	Produce(ctx context.Context, topic string, key, value []byte) error
	// Close releases the resources used by the Producer.
	Close() error
}

// Consumer is the interface for receiving messages from Kafka.
type Consumer interface {
	// Consume continuously reads messages and calls the provided handler for each message.
	Consume(ctx context.Context, handler func(topic string, key, value []byte) error) error
	// Close releases the resources used by the Consumer.
	Close() error
}

// Message represents a Kafka message with additional metadata.
type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Partition int
	Offset    int64
	Timestamp time.Time
}
