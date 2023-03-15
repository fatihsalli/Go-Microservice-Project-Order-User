package kafka

import (
	"OrderUserProject/internal/models"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
)

func ListenFromKafka(topic string) {
	// TODO: brokersUrl have to come config file
	// Kafka broker address
	brokersUrl := []string{"localhost:9092"}

	// Kafka configuration
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

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("Failed to create partition consumer: %v ", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Print(err)
		}
	}()

	var orderList []models.Order
	var order models.Order

	for msg := range partitionConsumer.Messages() {
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Print(err)
		}

		orderList = append(orderList, order)
		log.Printf("Received order: %+v\n", order)
	}
}
