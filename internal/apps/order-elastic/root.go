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
	return &OrderSyncService{
		Service:  service,
		Consumer: consumer,
		Config:   config,
		Producer: producer,
	}
}

func (r OrderSyncService) StartPushOrder() error {
	log.Info("Order sync service started")
	err := r.Consumer.SubscribeToTopics([]string{r.Config.Elasticsearch.TopicName["OrderID"]})
	if err != nil {
		log.Errorf("Kafka connection failed. | Error: %v\n", err)
	}
	for {
		ordersID := make([]string, 0)
		fromTopics, err := r.Consumer.ConsumeFromTopics(3, 10, 2)
		if err != nil {
			log.Info("An error when consume from topic...")
		}

		for _, message := range fromTopics {
			ordersID = append(ordersID, string(message.Value))
		}

		ordersModel, err := r.Service.GetOrderWithHttpClient(ordersID)
		if err != nil {
			log.Errorf("An error:%v", err)
		} else {
			r.Consumer.AckLastMessage()
		}

		for _, orderForPush := range ordersModel {
			// => SEND MESSAGE (Order Model)
			orderJSON, err := json.Marshal(orderForPush)
			if err != nil {
				log.Errorf("Error marshalling order:", err)
			}

			err = r.Producer.SendToKafkaWithMessage(orderJSON)
			if err != nil {
				log.Errorf("Something went wrong: %v", err)
			}
		}
	}
}

func (r OrderSyncService) StartConsumeOrder() error {
	log.Info("Order start to consume and save on elasticsearch!")
	err := r.Consumer.SubscribeToTopics([]string{r.Config.Elasticsearch.TopicName["OrderModel"]})
	if err != nil {
		log.Errorf("Kafka connection failed. | Error: %v\n", err)
	}

	for {
		fromTopics, err := r.Consumer.ConsumeFromTopics(3, 10, 2)
		if err != nil {
			log.Info("An error when consume from topic...")
		}

		for _, message := range fromTopics {
			var orderResponse OrderResponse
			jsonErr := json.Unmarshal(message.Value, &orderResponse)
			if jsonErr == nil {
				err = r.Service.SaveOrderToElasticsearch(orderResponse)
				if err != nil {
					log.Errorf("Something went wrong: %v", err)
				}
			} else {
				log.Errorf(jsonErr.Error())
			}
		}
	}
}
