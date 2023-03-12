package elasticsearch

import (
	"OrderUserProject/internal/models"
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"log"
)

func SaveElastic(topic string) {
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

	// Connect ElasticSearch client
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating ElasticSearch client: %s", err)
	}

	// Kafka listen from topic
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalln(err)
	}
	initialOffset := sarama.OffsetNewest

	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			log.Print(err)
		}

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				// Mesajı REST API'ye gönderme
				log.Printf("Message topic:%s partition:%d offset:%d value:%s\n", message.Topic, message.Partition, message.Offset, message.Value)

				// Save to ElasticSearch
				var order models.Order
				if err := json.Unmarshal(message.Value, &order); err != nil {
					log.Printf("Error unmarshalling message: %s", err)
				} else {
					req := esapi.IndexRequest{
						Index:      "orders",
						DocumentID: order.ID,
						Body:       esutil.NewJSONReader(order),
						Refresh:    "true",
					}

					res, err := req.Do(context.Background(), es)
					if err != nil {
						log.Printf("Error indexing document: %s", err)
					} else {
						defer res.Body.Close()
					}
				}
			}
		}(pc)
	}
}
