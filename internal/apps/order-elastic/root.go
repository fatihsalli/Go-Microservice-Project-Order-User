package order_elastic

import (
	"OrderUserProject/pkg/kafka"
	"github.com/sirupsen/logrus"
)

type OrderSyncService struct {
	service  *OrderElasticService
	consumer *kafka.KafkaConsumer
	logger   *logrus.Logger
}

func NewProductSyncService(service *OrderElasticService, consumer *kafka.KafkaConsumer, logger *logrus.Logger) *OrderSyncService {
	return &OrderSyncService{
		service:  service,
		consumer: consumer,
		logger:   logger}
}
func (r OrderSyncService) Start() {

}
