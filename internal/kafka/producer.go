package kafka

import (
	"github.com/Shopify/sarama"
	"log"
)

func SendToKafka(topic string, message []byte) error {
	// Kafka broker address
	brokersUrl := []string{"localhost:9092"}

	// Sarama set-up
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	// Connect Kafka
	producer, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Send message to Kafka
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
