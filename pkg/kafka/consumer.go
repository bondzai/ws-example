// pkg/kafka/consumer.go
package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type kafkaConsumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new instance of Consumer.
// brokers: a list of Kafka brokers
// groupID: the consumer group ID for subscription
// topics: the list of topics to subscribe to
func NewConsumer(brokers []string, groupID string, topics []string) Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: topics,
	})
	log.Println("connected to kafka customer group: " + groupID)
	return &kafkaConsumer{reader: reader}
}

// Consume reads messages from Kafka in a loop and calls the handler for each message.
func (c *kafkaConsumer) Consume(ctx context.Context, handler func(topic string, key, value []byte) error) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := handler(msg.Topic, msg.Key, msg.Value); err != nil {
			log.Printf("Handler error: %v", err)
		}
	}
}

// Close shuts down the reader to release resources.
func (c *kafkaConsumer) Close() error {
	return c.reader.Close()
}
