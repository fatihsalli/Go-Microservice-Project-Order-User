package order_elastic

import (
	order_api "OrderUserProject/internal/apps/order-api"
	"OrderUserProject/pkg/kafka"
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/labstack/gommon/log"
)

type OrderElasticService struct {
}

func NewOrderElasticService() *OrderElasticService {
	orderService := &OrderElasticService{}
	return orderService
}

type IOrderElasticService interface {
	ConsumeOrderDuplicate(topic string) (order_api.OrderResponse, error)
	SaveOrderToElasticsearch(order order_api.OrderResponse) error
}

func (b OrderElasticService) ConsumeOrderDuplicate(topic string) (order_api.OrderResponse, error) {
	// => RECEIVE MESSAGE
	result := kafka.ListenFromKafka(topic)
	var order order_api.OrderResponse

	err := json.Unmarshal(result, &order)
	if err != nil {
		return order_api.OrderResponse{}, err
	}

	return order, nil
}

func (b OrderElasticService) SaveOrderToElasticsearch(order order_api.OrderResponse) error {

	//config and client
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	esClient, err := elasticsearch.NewClient(cfg)

	// Build the request body.
	data, err := json.Marshal(order)
	if err != nil {
		log.Errorf("Error marshaling document: %s", err)
	}

	// Set up the request object.
	req := esapi.CreateRequest{
		Index:      "order-duplicate-V01",
		DocumentID: order.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		log.Errorf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), order.ID)
	}

	return nil
}
