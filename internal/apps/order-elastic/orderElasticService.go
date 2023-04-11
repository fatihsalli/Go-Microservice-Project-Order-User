package order_elastic

import (
	"OrderUserProject/internal/configs"
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/neko-neko/echo-logrus/v2/log"
)

type OrderElasticService struct {
}

func NewOrderElasticService() *OrderElasticService {
	orderElasticService := &OrderElasticService{}
	return orderElasticService
}

func (b *OrderElasticService) SaveOrderToElasticsearch(order OrderResponse, config configs.Config) error {
	// client with default config
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.Elasticsearch.Addresses["Address 1"],
		},
	}

	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Errorf("Error creating the client: ", err)
		return err
	}

	// Build the request body.
	data, err := json.Marshal(order)
	if err != nil {
		log.Errorf("Error marshaling document: %s", err)
		return err
	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      config.Elasticsearch.IndexName["OrderSave"],
		DocumentID: order.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		log.Errorf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Errorf("Error parsing the response body: %s", err)
			return err
		} else {
			// Print the error information.
			log.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}

func (b *OrderElasticService) DeleteOrderFromElasticsearch(orderID string, config configs.Config) error {
	// client with default config
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.Elasticsearch.Addresses["Address 1"],
		},
	}

	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Errorf("Error creating the client: ", err)
		return err
	}

	// Create request object
	req := esapi.DeleteRequest{
		Index:      config.Elasticsearch.IndexName["OrderSave"],
		DocumentID: orderID,
		Refresh:    "true",
	}

	// Execute the request
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Errorf("Error parsing the response body: %s", err)
			return err
		} else {
			// Print the error information.
			log.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}
