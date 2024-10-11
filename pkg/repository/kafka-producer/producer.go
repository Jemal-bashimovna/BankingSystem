package kafkaproducer

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func NewProducer(brokers, topic string) *Producer {
	config := &kafka.ConfigMap{"bootstrap.servers": brokers}
	p, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}

	return &Producer{producer: p, topic: topic}
}

func (p *Producer) SendMessage(message []byte) error {
	deliveryChan := make(chan kafka.Event)

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan // Wait for delivery report
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}

func (p *Producer) Close() {
	p.producer.Flush(15 * 1000)
	p.producer.Close()
}
