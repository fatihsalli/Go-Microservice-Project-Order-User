package roots

import (
	"OrderUserProject/internal/apps/order-elastic"
	"OrderUserProject/internal/configs"
	kafkaPackage "OrderUserProject/pkg/kafka"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type OrderEventRoot struct {
	ServiceEvent   *order_elastic.OrderEventService
	ServiceElastic *order_elastic.OrderElasticService
	Consumer       *kafkaPackage.ConsumerKafka
	Producer       *kafkaPackage.ProducerKafka
	Config         *configs.Config
	Logger         *logrus.Logger
}

func NewOrderEventRoot(serviceEvent *order_elastic.OrderEventService, serviceElastic *order_elastic.OrderElasticService, consumer *kafkaPackage.ConsumerKafka, producer *kafkaPackage.ProducerKafka, config *configs.Config, logger *logrus.Logger) *OrderEventRoot {
	return &OrderEventRoot{
		ServiceEvent:   serviceEvent,
		ServiceElastic: serviceElastic,
		Consumer:       consumer,
		Config:         config,
		Producer:       producer,
		Logger:         logger,
	}
}

// StartGetOrderAndPushOrder => Get message from Kafka to consume OrderID, get order with http.client and push order with Kafka
func (o *OrderEventRoot) StartGetOrderAndPushOrder() error {
	o.Logger.Info("OrderSyncService starting for consume 'OrderID'.")
	err := o.Consumer.SubscribeToTopics([]string{o.Config.Kafka.TopicName["OrderID"]})
	if err != nil {
		o.Logger.Errorf("Kafka connection failed. | Error: %v\n", err)
	}
	for {
		ordersID := make([]string, 0)
		fromTopics, err := o.Consumer.ConsumeFromTopics(1, 5, 2)
		if err != nil {
			o.Logger.Errorf("An error when consume from topic. | Error: %v\n", err)
		}

		for _, message := range fromTopics {
			o.Logger.Infof("Message received from kafka: %v\n", string(message.Value))

			var orderResponse order_elastic.OrderResponseForElastic
			jsonErr := json.Unmarshal(message.Value, &orderResponse)
			if jsonErr != nil {
				o.Logger.Errorf(jsonErr.Error())
			}

			switch orderResponse.Status {
			case "Created", "Updated":
				ordersID = append(ordersID, orderResponse.OrderID)
			case "Deleted":
				if err := o.ServiceElastic.DeleteOrderFromElasticsearch(orderResponse.OrderID, *o.Config); err != nil {
					o.Logger.Errorf("An error deleting order from es. | Error: %v\n", err)
				}
				o.Logger.Infof("Order (ID:%v) successfully deleted from es.", orderResponse.OrderID)
			default:
				o.Logger.Errorf("Unknown order response status. | Error: %v\n", orderResponse.Status)
			}
		}

		if len(ordersID) > 0 {
			ordersModel, err := o.ServiceEvent.GetOrderWithHttpClient(ordersID, o.Config.HttpClient.OrderAPI)
			if err != nil || ordersModel == nil {
				o.Logger.Errorf("Orders cannot find. | Error: %v\n", err)
			} else {
				o.Consumer.AckLastMessage()
			}

			for _, orderForPush := range ordersModel {
				// => SEND MESSAGE (Order Model)
				orderJSON, err := json.Marshal(orderForPush)
				if err != nil {
					o.Logger.Errorf("An error when convert from json. | Error: %v\n", err)
				}

				err = o.Producer.SendToKafkaWithMessage(orderJSON, o.Config.Kafka.TopicName["OrderModel"])
				if err != nil {
					o.Logger.Errorf("An error when send a message... | Error: %v\n", err)
				} else {
					o.Logger.Infof("Order successfully pushed with id: %v", orderForPush.ID)
				}
			}
		}
	}
}
