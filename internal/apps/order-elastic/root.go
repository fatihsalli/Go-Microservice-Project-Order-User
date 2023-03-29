package order_elastic

import (
	kafkaConsumer "OrderUserProject/pkg/kafka"
	"github.com/neko-neko/echo-logrus/v2/log"
)

type OrderSyncService struct {
	service  *OrderElasticService
	consumer *kafkaConsumer.ConsumerKafka
}

func NewOrderSyncService(service *OrderElasticService, consumer *kafkaConsumer.ConsumerKafka) *OrderSyncService {
	return &OrderSyncService{service: service, consumer: consumer}
}

func (r OrderSyncService) Start(topic string) error {
	log.Info("Order sync service started")
	err := r.consumer.SubscribeToTopics([]string{topic})
	if err != nil {
		log.Errorf("Kafka connection failed. | Error: %v\n", err)
	}

	return nil
}
