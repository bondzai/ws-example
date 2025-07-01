// pkg/kafka/producer.go
package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

// NewProducer creates a new instance of Producer.
// brokers: a list of Kafka brokers (e.g., []string{"localhost:9092"})
// defaultTopic: the default topic to use (can be overridden in Produce)
func NewProducer(brokers []string, defaultTopic string) Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        defaultTopic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 100 * time.Millisecond,
	})
	return &kafkaProducer{writer: writer}
}

// Produce sends a message to the specified topic. It can override the default topic.
func (p *kafkaProducer) Produce(ctx context.Context, topic string, key, value []byte) error {
	p.writer.Topic = topic
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}

// Close shuts down the writer to release resources.
func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}
