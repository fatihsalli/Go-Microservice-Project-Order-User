package order_elastic

import (
	"OrderUserProject/internal/models"
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
	ConsumeOrderDuplicate(topic string) (models.Order, error)
	SaveOrderToElasticsearch(order models.Order) error
}

func (b OrderElasticService) ConsumeOrderDuplicate(topic string) (models.Order, error) {
	// => RECEIVE MESSAGE
	result := kafka.ListenFromKafka(topic)
	var orderResponse OrderResponse

	err := json.Unmarshal(result, &orderResponse)
	if err != nil {
		return models.Order{}, err
	}

	var order models.Order

	// we can use automapper, but it will cause performance loss.
	order.ID = orderResponse.ID
	order.UserId = orderResponse.UserId
	order.Status = orderResponse.Status
	order.Address.Address = orderResponse.Address.Address
	order.Address.City = orderResponse.Address.City
	order.Address.District = orderResponse.Address.District
	order.Address.Type = orderResponse.Address.Type
	order.Address.Default = orderResponse.Address.Default
	order.InvoiceAddress.Address = orderResponse.InvoiceAddress.Address
	order.InvoiceAddress.City = orderResponse.InvoiceAddress.City
	order.InvoiceAddress.District = orderResponse.InvoiceAddress.District
	order.InvoiceAddress.Type = orderResponse.InvoiceAddress.Type
	order.InvoiceAddress.Default = orderResponse.InvoiceAddress.Default
	order.Total = orderResponse.Total
	order.Product = orderResponse.Product
	order.CreatedAt = orderResponse.CreatedAt
	order.UpdatedAt = orderResponse.UpdatedAt

	return order, nil
}

func (b OrderElasticService) SaveOrderToElasticsearch(order models.Order) error {
	// client with default config => http://localhost:9200
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
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

	// TODO : versiyonlama araştırılacak
	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      "order_duplicate_v01",
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
