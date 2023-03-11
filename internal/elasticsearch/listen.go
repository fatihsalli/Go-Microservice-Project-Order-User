package elasticsearch

import (
	"OrderUserProject/internal/kafka"
	"OrderUserProject/internal/models"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

func ListenTopic() {
	topic := "order-create"
	consumer, err := kafka.ConnectConsumer([]string{"localhost:9092"})

	if err != nil {
		panic(err)
	}

	ordersConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)

	// Start consuming messages
	for msg := range ordersConsumer.Messages() {
		order := &models.Order{}
		if err := json.Unmarshal(msg.Value, order); err != nil {
			fmt.Println("Failed to unmarshal order: ", err)
		} else {
			fmt.Printf("Received order: %+v\n", order)
		}
	}

	return
}
