package kafka

import (
	"github.com/Shopify/sarama"
	"log"
)

func ListenFromKafka(topic string) {
	// Kafka broker address
	brokersUrl := []string{"localhost:9092"}

	// Sarama set-up
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Connect Kafka
	consumer, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Kafka listen from topic
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalln(err)
	}
	initialOffset := sarama.OffsetOldest
	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			log.Print(err)
		}

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				// Mesajı REST API'ye gönderme
				log.Printf("Message topic:%s partition:%d offset:%d value:%s\n", message.Topic, message.Partition, message.Offset, message.Value)
				// işlem yap
			}
		}(pc)
	}
}
