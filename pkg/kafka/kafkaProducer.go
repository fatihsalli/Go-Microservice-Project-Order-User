package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/labstack/gommon/log"
)

// SendToKafka take a topic name and message with format of []byte
func SendToKafka(topic string, msg string) {
	// Kafka configuration
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Kafka broker address
	brokers := []string{"localhost:9092"}

	// Kafka producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Kafka consumer
	consumer, err := sarama.NewConsumerGroup(brokers, "my-group", config)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Kafka consumer starting with goroutine
	go func() {
		for {
			topics := []string{topic}
			handler := &orderHandler{producer: producer}

			err := consumer.Consume(context.Background(), topics, handler)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	// Send a message with Kafka producer
	message := &sarama.ProducerMessage{
		Topic: "my-topic",
		Value: sarama.StringEncoder(msg),
	}
	_, _, err = producer.SendMessage(message)
	if err != nil {
		log.Fatal(err)
	}
}

type orderHandler struct {
	producer sarama.SyncProducer
}

func (h *orderHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *orderHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *orderHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		orderID := string(message.Value)
		// OrderID'yi başka bir fonksiyona göndermek için burada işlem yapabilirsiniz
		log.Printf("Received orderID: %s\n", orderID)

		session.MarkMessage(message, "")
	}
	return nil
}
