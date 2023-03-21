package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
)

// Kafka config
func newKafkaConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}
}

// SendToKafka take a topic name and message with format of []byte
func SendToKafka(topic string, msg string) {
	// Kafka configuration
	config := newKafkaConfig()

	// Producer
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Errorf("error creating producer: %v", err)
	}

	// To prepare message
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(msg),
	}

	// Send to message
	err = producer.Produce(message, nil)
	if err != nil {
		log.Errorf("error producing message: %v", err)
	}

	// Close producer
	producer.Flush(15 * 1000)
}
