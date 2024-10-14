package deps

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
}

func NewConsumer(brokers, groupId string, topics []string) *Consumer {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	}
	c, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}

	c.SubscribeTopics(topics, nil)

	return &Consumer{consumer: c}
}

func (c *Consumer) PollMessage() ([]byte, error) {
	msg, err := c.consumer.ReadMessage(-1)
	if err != nil {
		return nil, err
	}
	return msg.Value, nil
}

func (c *Consumer) Close() {
	c.consumer.Close()
}
