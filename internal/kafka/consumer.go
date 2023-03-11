package kafka

import (
	"github.com/Shopify/sarama"
	"log"
)

func ConnectConsumer() (sarama.Consumer, error) {
	// Kafka broker address
	brokerList := []string{"localhost:9092"}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Kafka consumer
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Printf("Error creating Kafka consumer: %s", err.Error())
	}
	defer consumer.Close()

	return consumer, nil
}
