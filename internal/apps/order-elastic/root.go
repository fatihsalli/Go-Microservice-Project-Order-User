package order_elastic

import (
	"OrderUserProject/internal/configs"
	kafkaPackage "OrderUserProject/pkg/kafka"
	"encoding/json"
	"github.com/labstack/gommon/log"
)

type OrderSyncService struct {
	Service  *OrderElasticService
	Consumer *kafkaPackage.ConsumerKafka
	Producer *kafkaPackage.ProducerKafka
	Config   *configs.Config
}

func NewOrderSyncService(service *OrderElasticService, consumer *kafkaPackage.ConsumerKafka, producer *kafkaPackage.ProducerKafka, config *configs.Config) *OrderSyncService {
	return &OrderSyncService{Service: service, Consumer: consumer, Config: config, Producer: producer}
}

func (r OrderSyncService) Start(topic string) error {
	result, err := r.Consumer.ListenFromKafkaWithoutTopic(r.Config.Elasticsearch.TopicName["OrderID"])
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
		return err
	}

	orderResponse, err := r.Service.GetOrderWithHttpClient(string(result))
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
		return err
	}

	// => SEND MESSAGE (Order Model)
	orderJSON, err := json.Marshal(orderResponse)
	if err != nil {
		log.Errorf("Error marshalling order:", err)
		return err
	}

	err = r.Producer.SendToKafkaWithMessage(orderJSON)
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
		return err
	}

	order, err := r.Service.ConsumeOrderDuplicate()
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
		return err
	}

	err = r.Service.SaveOrderToElasticsearch(order)
	if err != nil {
		log.Errorf("Something went wrong: %v", err)
		return err
	}

	return nil
}
