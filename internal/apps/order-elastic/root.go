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

// StartGetOrderAndPushOrder => Get message from Kafka to consume OrderID, get order with http.client and push order with Kafka
func (r *OrderSyncService) StartGetOrderAndPushOrder() error {
	r.Logger.Info("OrderSyncService starting for consume 'OrderID'.")
	err := r.Consumer.SubscribeToTopics([]string{r.Config.Kafka.TopicName["OrderID"]})
	if err != nil {
		r.Logger.Errorf("Kafka connection failed. | Error: %v\n", err)
	}
	for {
		ordersID := make([]string, 0)
		fromTopics, err := r.Consumer.ConsumeFromTopics(1, 5, 2)
		if err != nil {
			r.Logger.Errorf("An error when consume from topic. | Error: %v\n", err)
		}

		for _, message := range fromTopics {
			r.Logger.Infof("Message received from kafka: %v\n", string(message.Value))

			var orderResponse OrderResponseForElastic
			jsonErr := json.Unmarshal(message.Value, &orderResponse)
			if jsonErr != nil {
				r.Logger.Errorf(jsonErr.Error())
			}

			switch orderResponse.Status {
			case "Created", "Updated":
				ordersID = append(ordersID, orderResponse.OrderID)
			case "Deleted":
				if err := r.Service.DeleteOrderFromElasticsearch(orderResponse.OrderID, *r.Config); err != nil {
					r.Logger.Errorf("An error deleting order from elasticsearch. | Error: %v\n", err)
				}
				r.Logger.Infof("Order (ID:%v) successfully deleted from elasticsearch.", orderResponse.OrderID)
			default:
				r.Logger.Errorf("Unknown order response status. | Error: %v\n", orderResponse.Status)
			}
		}

		if len(ordersID) > 0 {
			ordersModel, err := r.Service.GetOrderWithHttpClient(ordersID)
			if err != nil || ordersModel == nil {
				r.Logger.Errorf("Orders cannot find. | Error: %v\n", err)
			} else {
				r.Consumer.AckLastMessage()
			}

			for _, orderForPush := range ordersModel {
				// => SEND MESSAGE (Order Model)
				orderJSON, err := json.Marshal(orderForPush)
				if err != nil {
					r.Logger.Errorf("An error when convert from json. | Error: %v\n", err)
				}

				err = r.Producer.SendToKafkaWithMessage(orderJSON, r.Config.Kafka.TopicName["OrderModel"])
				if err != nil {
					r.Logger.Errorf("An error when send a message... | Error: %v\n", err)
				} else {
					r.Logger.Infof("Order successfully pushed with id: %v", orderForPush.ID)
				}
			}
		}
	}
}

// StartConsumeAndSaveOrder => Get message from Kafka to consume OrderModel and save/update on elasticsearch
func (r *OrderSyncService) StartConsumeAndSaveOrder() error {
	r.Logger.Info("OrderSyncService starting to consume 'OrderModel'.")
	err := r.Consumer.SubscribeToTopics([]string{r.Config.Kafka.TopicName["OrderModel"]})
	if err != nil {
		r.Logger.Errorf("Kafka connection failed. | Error: %v\n", err)
	}

	for {
		fromTopics, err := r.Consumer.ConsumeFromTopics(1, 5, 2)
		if err != nil {
			r.Logger.Errorf("An error when consume from topic. | Error: %v\n", err)
		}

		for _, message := range fromTopics {
			var orderResponse OrderResponse
			jsonErr := json.Unmarshal(message.Value, &orderResponse)
			if jsonErr == nil {
				err = r.Service.SaveOrderToElasticsearch(orderResponse, *r.Config)
				if err != nil {
					r.Logger.Errorf("Order cannot save on elasticsearch. | Error: %v\n", err)
				} else {
					r.Logger.Infof("Order (ID:%v) saved on elasticsearch.", orderResponse.ID)
				}
			} else {
				r.Logger.Errorf("An error when convert to json. | Error: %v\n", jsonErr.Error())
			}
		}
	}
}
