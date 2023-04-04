package order_elastic

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type OrderEventService struct {
	Logger *logrus.Logger
}

func NewOrderEventService(logger *logrus.Logger) *OrderEventService {
	orderEventService := &OrderEventService{Logger: logger}
	return orderEventService
}

func (o *OrderEventService) GetOrderWithHttpClient(ordersID []string) ([]OrderResponse, error) {

	var orders []OrderResponse

	for _, orderID := range ordersID {
		// => HTTP.CLIENT FIND ORDER
		// Create a new HTTP client with a timeout
		client := http.Client{
			Timeout: time.Second * 20,
		}

		// Send a GET request to the Order service to retrieve order information
		respOrder, err := client.Get("http://localhost:8011/api/orders" + "/" + orderID)
		if err != nil || respOrder.StatusCode != http.StatusOK {
			o.Logger.Errorf("Order with id {%v} cannot find!", orderID)
			return []OrderResponse{}, err
		}

		// Read the response body
		respOrderBody, err := io.ReadAll(respOrder.Body)
		if err != nil {
			o.Logger.Errorf("StatusInternalServerError: %v", err.Error())
			return []OrderResponse{}, err
		}

		// Unmarshal the response body into an Order struct
		var orderResponse OrderResponse
		err = json.Unmarshal(respOrderBody, &orderResponse)
		if err != nil {
			o.Logger.Errorf("StatusInternalServerError: %v", err.Error())
			return []OrderResponse{}, err
		}

		orders = append(orders, orderResponse)
	}

	return orders, nil
}
