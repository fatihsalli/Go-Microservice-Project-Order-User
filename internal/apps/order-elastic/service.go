package order_elastic

import (
	"OrderUserProject/internal/models"
	"OrderUserProject/pkg/kafka"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	ConsumeOrderDuplicate(topic string) (models.Order, error)
	SaveOrderToElasticsearch(order models.Order) error
	GetOrderFromElasticsearch(orderID string) (models.Order, error)
}

var (
	esClient *elasticsearch.Client
)

func init() {
	var err error
	esClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Errorf("Error creating the client: %s", err.Error())
	}
}

func (b OrderElasticService) ConsumeOrderDuplicate(topic string) (models.Order, error) {
	// => RECEIVE MESSAGE
	result := kafka.ListenFromKafka(topic)
	var order models.Order

	err := json.Unmarshal(result, &order)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (b OrderElasticService) SaveOrderToElasticsearch(order models.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("error marshalling order to JSON: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      "orders",
		DocumentID: order.ID,
		Body:       bytes.NewReader(orderJSON),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.Status())
	}

	return nil
}

func (b OrderElasticService) GetOrderFromElasticsearch(orderID string) (models.Order, error) {
	req := esapi.GetRequest{
		Index:      "orders",
		DocumentID: orderID,
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return models.Order{}, fmt.Errorf("error getting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return models.Order{}, fmt.Errorf("error getting document: %s", res.Status())
	}

	var order models.Order
	if err := json.NewDecoder(res.Body).Decode(&order); err != nil {
		return models.Order{}, fmt.Errorf("error decoding document: %w", err)
	}

	return order, nil
}
