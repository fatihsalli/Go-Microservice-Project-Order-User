package order_elastic

import (
	"OrderUserProject/internal/configs"
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"time"
)

type OrderElasticService struct {
}

func NewOrderElasticService() *OrderElasticService {
	orderElasticService := &OrderElasticService{}
	return orderElasticService
}

var ClientBaseUrl = map[string]string{
	"order":         "http://localhost:8011/api/orders",
	"user":          "http://localhost:8012/api/users",
	"order-elastic": "http://localhost:8013/api/orders-elastic",
}

type IOrderElasticService interface {
	GetOrderWithHttpClient(orderID string) (OrderResponse, error)
	SaveOrderToElasticsearch(order OrderResponse) error
}

func (b *OrderElasticService) GetOrderWithHttpClient(ordersID []string) ([]OrderResponse, error) {

	var orders []OrderResponse

	for _, orderID := range ordersID {
		// => HTTP.CLIENT FIND ORDER
		// Create a new HTTP client with a timeout
		client := http.Client{
			Timeout: time.Second * 20,
		}

		// Send a GET request to the Order service to retrieve order information
		respOrder, err := client.Get(ClientBaseUrl["order"] + "/" + orderID)
		if err != nil || respOrder.StatusCode != http.StatusOK {
			log.Errorf("Order with id {%v} cannot find!", orderID)
			return []OrderResponse{}, err
		}

		// Read the response body
		respOrderBody, err := io.ReadAll(respOrder.Body)
		if err != nil {
			log.Errorf("StatusInternalServerError: %v", err.Error())
			return []OrderResponse{}, err
		}

		// Unmarshal the response body into an Order struct
		var orderResponse OrderResponse
		err = json.Unmarshal(respOrderBody, &orderResponse)
		if err != nil {
			log.Errorf("StatusInternalServerError: %v", err.Error())
			return []OrderResponse{}, err
		}

		orders = append(orders, orderResponse)
	}

	return orders, nil
}

func (b *OrderElasticService) SaveOrderToElasticsearch(order OrderResponse, config configs.Config) error {
	// client with default config => http://localhost:9200
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

	// TODO : versiyonlama araştırılacak
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
