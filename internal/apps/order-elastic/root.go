package order_elastic

import (
	"OrderUserProject/internal/configs"
	kafkaPackage "OrderUserProject/pkg/kafka"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type OrderSyncService struct {
	Service  *OrderElasticService
	Consumer *kafkaPackage.ConsumerKafka
	Producer *kafkaPackage.ProducerKafka
	Config   *configs.Config
	Logger   *logrus.Logger
}

func NewOrderSyncService(service *OrderElasticService, consumer *kafkaPackage.ConsumerKafka, producer *kafkaPackage.ProducerKafka, config *configs.Config, logger *logrus.Logger) *OrderSyncService {
	return &OrderSyncService{
		Service:  service,
		Consumer: consumer,
		Config:   config,
		Producer: producer,
		Logger:   logger,
	}
}

// StartPushOrder => Get message from Kafka to consume OrderID, get order with http.client and push order with Kafka
func (r *OrderSyncService) StartPushOrder() error {
	r.Logger.Info("Order sync service start to get OrderID and push Order models!")
	err := r.Consumer.SubscribeToTopics([]string{r.Config.Kafka.TopicName["OrderID"]})
	if err != nil {
		r.Logger.Errorf("Kafka connection failed. | Error: %v\n", err)
	}
	for {
		ordersID := make([]string, 0)
		fromTopics, err := r.Consumer.ConsumeFromTopics(1, 5, 2)
		if err != nil {
			r.Logger.Errorf("An error when consume from topic...:%v", err)
		}

		for _, message := range fromTopics {
			r.Logger.Infof("Messaige received: %v", string(message.Value))
			ordersID = append(ordersID, string(message.Value))
		}

		ordersModel, err := r.Service.GetOrderWithHttpClient(ordersID)
		if err != nil {
			r.Logger.Errorf("An error:%v", err)
		} else {
			r.Consumer.AckLastMessage()
		}

		for _, orderForPush := range ordersModel {
			// => SEND MESSAGE (Order Model)
			orderJSON, err := json.Marshal(orderForPush)
			if err != nil {
				r.Logger.Errorf("Error marshalling order: %v", err)
			}

			err = r.Producer.SendToKafkaWithMessage(orderJSON, r.Config.Kafka.TopicName["OrderModel"])
			if err != nil {
				r.Logger.Errorf("Something went wrong: %v", err)
			} else {
				r.Logger.Infof("Order pushed with id: %v", orderForPush.ID)
			}
		}
	}
}

// StartConsumeOrder => Get message from Kafka to consume OrderModel and save on elasticsearch
func (r *OrderSyncService) StartConsumeOrder() error {
	r.Logger.Info("Order sync service start to consume Order Model and save on elasticsearch!")
	err := r.Consumer.SubscribeToTopics([]string{r.Config.Kafka.TopicName["OrderModel"]})
	if err != nil {
		r.Logger.Errorf("Kafka connection failed. | Error: %v\n", err)
	}

	for {
		fromTopics, err := r.Consumer.ConsumeFromTopics(1, 5, 2)
		if err != nil {
			r.Logger.Errorf("An error when consume from topic...: %v", err)
		}

		for _, message := range fromTopics {
			var orderResponse OrderResponse
			jsonErr := json.Unmarshal(message.Value, &orderResponse)
			if jsonErr == nil {
				err = r.Service.SaveOrderToElasticsearch(orderResponse, *r.Config)
				if err != nil {
					r.Logger.Errorf("Something went wrong: %v", err)
				} else {
					r.Logger.Infof("Order (%v) saved on elasticsearch", orderResponse.ID)
				}
			} else {
				r.Logger.Errorf(jsonErr.Error())
			}
		}
	}
}
