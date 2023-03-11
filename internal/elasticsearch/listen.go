package elasticsearch

import (
	"OrderUserProject/internal/kafka"
	"fmt"
	"github.com/Shopify/sarama"
)

func ListenTopic() {
	topic := "order-create"
	worker, err := kafka.ConnectConsumer([]string{"localhost:9092"})

	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)

	fmt.Println(consumer.Messages())
}
