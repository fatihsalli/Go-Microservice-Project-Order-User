package kafka

import (
	"github.com/Shopify/sarama"
	"log"
)

func ListenFromKafka(topic string) sarama.Consumer {
	// Kafka broker address
	brokersUrl := []string{"localhost:9092"}

	// Sarama set-up
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Connect Kafka
	consumer, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Print(err)
		}
	}()

	return consumer
}
