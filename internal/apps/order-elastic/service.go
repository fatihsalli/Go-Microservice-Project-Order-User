package order_elastic

import (
	"OrderUserProject/internal/models"
	"OrderUserProject/pkg/kafka"
	"encoding/json"
)

type OrderElasticService struct {
}

func NewOrderElasticService() *OrderElasticService {
	orderService := &OrderElasticService{}
	return orderService
}

type IOrderElasticService interface {
}

func CreateOrderDuplicate() error {
	// => RECEIVE MESSAGE
	// create topic name
	topic := "orderDuplicate-created-v01"

	result := kafka.ListenFromKafka(topic)
	var order models.Order

	err := json.Unmarshal(result, &order)
	if err != nil {
		return err
	}

	return nil
}
