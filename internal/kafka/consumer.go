package kafka

import (
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/repository"
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

	var order models.Order

	for msg := range partitionConsumer.Messages() {
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Print(err)
		}

		// elastic.SaveOrderElastic(order)

		log.Printf("Received order: %+v\n", order)
	}
}

// SaveOrder for test to consume event and write on MongoDB
func SaveOrder(order models.Order) {
	//for test
	config := configs.GetConfig("test")
	mongoOrderCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection("Orders-event")
	OrderRepository := repository.NewOrderRepository(mongoOrderCollection)

	result, err := OrderRepository.Insert(order)
	if result == false || err != nil {
		log.Printf("Cannot create order event in MongoDB %v", order.ID)
	}
}
