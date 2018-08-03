package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer reads from a topic
type Consumer struct {
	kc *kafka.Consumer
}

// New creates a new Consumer.
func New(bootstrapServers, consumerGroup, offsetResetType, topic string) (c *Consumer) {

	kc, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          consumerGroup,
		"auto.offset.reset": offsetResetType,
	})
	if err != nil {
		log.Fatal(err)
	}

	kc.SubscribeTopics([]string{topic, "^aRegex.*[Tt]opic"}, nil)

	c = &Consumer{kc: kc}

	return c
}

// Consume consumes messages from a topic
func (c *Consumer) Consume() (msg *kafka.Message, err error) {
	msg, err = c.kc.ReadMessage(-1)
	return msg, err
}

// Close closes the kafka connection
func (c *Consumer) Close() {
	c.kc.Close()
}
