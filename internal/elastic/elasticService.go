package elastic

import (
	"OrderUserProject/internal/models"
	"context"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"github.com/olivere/elastic/v7"
)

func SaveOrderElastic(order models.Order) {
	// Elasticsearch configuration
	esUrl := "http://localhost:9200"
	esIndex := "orders-duplicate"

	// Create Elasticsearch client
	client, err := elastic.NewClient(
		elastic.SetURL(esUrl),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Printf("Failed to create Elasticsearch client: %v ", err)
		return
	}

	// Save order to Elasticsearch
	_, err = client.Index().
		Index(esIndex).
		BodyJson(order).
		Do(context.Background())
	if err != nil {
		log.Printf("Failed to save order to Elasticsearch: %v ", err)
		return
	}

	log.Printf("Saved order to Elasticsearch: %+v\n", order)
}

func ReadFromElastic() {
	// Create Elasticsearch client
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	searchResult, err := client.Search().
		Index("orders-duplicate").
		Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var orders []models.Order
	for _, hit := range searchResult.Hits.Hits {
		var order models.Order
		err := json.Unmarshal(hit.Source, &order)
		if err != nil {
			log.Printf("Failed to unmarshal order: %v", err)
		} else {
			orders = append(orders, order)
		}
	}

	log.Printf("Elastic work successfully! Found %d orders", len(orders))
}
