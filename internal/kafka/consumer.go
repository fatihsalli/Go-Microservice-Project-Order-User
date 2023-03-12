package kafka

import (
	"OrderUserProject/internal/models"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
)

func ListenFromKafka(topic string) []models.Order {
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

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to create partition consumer: %v ", err)
	}
	defer partitionConsumer.Close()

	var orderList []models.Order
	var order models.Order

	// Kafka listen from topic
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Print(err)
	}
	initialOffset := sarama.OffsetNewest
	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			log.Print(err)
		}

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				log.Printf("Message topic:%s partition:%d offset:%d value:%s\n", message.Topic, message.Partition, message.Offset, message.Value)

				if err := json.Unmarshal(message.Value, &order); err != nil {
					log.Printf("Error unmarshalling message: %s", err)
				}
				orderList = append(orderList, order)
			}
		}(pc)
	}

	return orderList

}
