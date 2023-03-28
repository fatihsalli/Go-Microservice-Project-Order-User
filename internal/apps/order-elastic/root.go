package order_elastic

import (
	"OrderUserProject/pkg/kafka"
	"github.com/sirupsen/logrus"
)

type ProductSyncService struct {
	service  *OrderElasticService
	consumer *kafka.KafkaConsumer
	logger   *logrus.Logger
}

func NewProductSyncService(service *OrderElasticService, consumer *kafka.KafkaConsumer, logger *logrus.Logger) *ProductSyncService {
	return &ProductSyncService{
		service:  service,
		consumer: consumer,
		logger:   logger}
}
func (r ProductSyncService) Start() {

}
