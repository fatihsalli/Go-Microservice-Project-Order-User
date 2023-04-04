package roots

import (
	"OrderUserProject/internal/apps/order-elastic"
	"OrderUserProject/internal/configs"
	kafkaPackage "OrderUserProject/pkg/kafka"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type OrderElasticRoot struct {
	Service  *order_elastic.OrderElasticService
	Consumer *kafkaPackage.ConsumerKafka
	Producer *kafkaPackage.ProducerKafka
	Config   *configs.Config
	Logger   *logrus.Logger
}

func NewOrderElasticRoot(service *order_elastic.OrderElasticService, consumer *kafkaPackage.ConsumerKafka, producer *kafkaPackage.ProducerKafka, config *configs.Config, logger *logrus.Logger) *OrderElasticRoot {
	return &OrderElasticRoot{
		Service:  service,
		Consumer: consumer,
		Config:   config,
		Producer: producer,
		Logger:   logger,
	}
}

// StartConsumeAndSaveOrder => Get message from Kafka to consume OrderModel and save/update on es
func (o *OrderElasticRoot) StartConsumeAndSaveOrder() error {
	o.Logger.Info("OrderSyncService starting to consume 'OrderModel'.")
	err := o.Consumer.SubscribeToTopics([]string{o.Config.Kafka.TopicName["OrderModel"]})
	if err != nil {
		o.Logger.Errorf("Kafka connection failed. | Error: %v\n", err)
	}

	for {
		fromTopics, err := o.Consumer.ConsumeFromTopics(1, 5, 2)
		if err != nil {
			o.Logger.Errorf("An error when consume from topic. | Error: %v\n", err)
		}

		for _, message := range fromTopics {
			var orderResponse order_elastic.OrderResponse
			jsonErr := json.Unmarshal(message.Value, &orderResponse)
			if jsonErr == nil {
				err = o.Service.SaveOrderToElasticsearch(orderResponse, *o.Config)
				if err != nil {
					o.Logger.Errorf("Order cannot save on es. | Error: %v\n", err)
				} else {
					o.Logger.Infof("Order (ID:%v) saved on es.", orderResponse.ID)
				}
			} else {
				o.Logger.Errorf("An error when convert to json. | Error: %v\n", jsonErr.Error())
			}
		}
	}
}
