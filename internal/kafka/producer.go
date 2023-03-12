package kafka

import (
	"github.com/Shopify/sarama"
	"log"
)

// SendToKafka take a topic name and message with format of []byte
func SendToKafka(topic string, message []byte) error {
	// TODO: brokersUrl have to come config file
	// Kafka broker address
	brokersUrl := []string{"localhost:9092"}

	// Kafka configuration
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	// Connect Kafka
	producer, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Print(err)
		}
	}()

	// Send message to Kafka
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)

	return nil
}
