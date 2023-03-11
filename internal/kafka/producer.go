package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/labstack/gommon/log"
)

func ConnectProducer() (sarama.SyncProducer, error) {
	// Kafka broker address
	brokerList := []string{"localhost:9092"}

	// Kafka setup settings
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	// Kafka producer
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Printf("Error creating Kafka producer: %s", err.Error())
	}
	defer producer.Close()

	return producer, nil
}
