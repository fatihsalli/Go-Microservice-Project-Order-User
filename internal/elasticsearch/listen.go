package elasticsearch

import (
	"OrderUserProject/internal/models"
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"log"
)

func ListenAndCreateTopic() {
	esCfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Printf("Error creating Elasticsearch client: %s", err)
	}

	consumerCfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, consumerCfg)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %s", err)
	}

	topic := "order-create-elastic"
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Error getting partitions for topic %s: %s", topic, err)
	}

	ctx := context.Background()

	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error creating partition consumer for partition %d: %s", partition, err)
		}

		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()

			for msg := range pc.Messages() {
				var order models.Order
				if err := json.Unmarshal(msg.Value, &order); err != nil {
					log.Printf("Error unmarshalling order: %s", err)
					continue
				}

				// Elasticsearch'e kaydet
				req := esutil.IndexRequest{
					Index:      "orders",
					DocumentID: order.ID,
					Body:       esutil.NewJSONReader(order),
					Refresh:    "true",
				}

				res, err := req.Do(ctx, es)
				if err != nil {
					log.Printf("Error indexing order: %s", err)
					continue
				}
				defer res.Body.Close()
			}
		}(pc)
	}
}
