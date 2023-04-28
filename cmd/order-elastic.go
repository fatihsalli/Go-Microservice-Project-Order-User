package cmd

import (
	"OrderUserProject/internal/apps/order-elastic"
	"OrderUserProject/internal/apps/order-elastic/roots"
	"OrderUserProject/internal/configs"
	"OrderUserProject/pkg/kafka"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func StartOrderElastic() {
	// Logger instead of standard log we use 'logrus' package
	logger := logrus.StandardLogger()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	logger.Info("Logger enabled!!")

	// Get config
	config := configs.GetConfig("test")

	// Create OrderElasticRoot => Consume orderModel, save on elastic search
	orderElasticService := order_elastic.NewOrderElasticService()
	producerElastic := kafka.NewProducerKafka(config.Kafka.Address)
	consumerElastic := kafka.NewConsumerKafka(config.Kafka.Address)
	orderElasticRoot := roots.NewOrderElasticRoot(orderElasticService, consumerElastic, producerElastic, &config, logger)

	// Create OrderEventRoot => Consume orderID, get order model, delete order from elastic and push order model
	orderEventService := order_elastic.NewOrderEventService(logger)
	producerEvent := kafka.NewProducerKafka(config.Kafka.Address)
	consumerEvent := kafka.NewConsumerKafka(config.Kafka.Address)
	orderEventRoot := roots.NewOrderEventRoot(orderEventService, orderElasticService, consumerEvent, producerEvent, &config, logger)

	// Create OrderSyncService
	orderSyncService := roots.NewOrderSyncService(orderElasticRoot, orderEventRoot)

	logger.Info("Order Elastic Service is starting...")
	orderSyncService.Start()
}
