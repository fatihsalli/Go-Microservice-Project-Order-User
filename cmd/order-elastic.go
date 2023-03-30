package cmd

import (
	order_elastic "OrderUserProject/internal/apps/order-elastic"
	"OrderUserProject/internal/configs"
	kafkaPackage "OrderUserProject/pkg/kafka"
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

	// First orderSyncService: Consume orderID, get order model and push order model
	// Create service,producer and consumer for orderSyncService

	service1 := order_elastic.NewOrderElasticService()
	producer1 := kafkaPackage.NewProducerKafka(config.Kafka.Address)
	consumer1 := kafkaPackage.NewConsumerKafka()
	orderSyncService1 := order_elastic.NewOrderSyncService(service1, consumer1, producer1, &config, logger)

	// Second orderSyncService: Consume orderModel, save on elastic search
	// Create service,producer and consumer for orderSyncService
	service2 := order_elastic.NewOrderElasticService()
	producer2 := kafkaPackage.NewProducerKafka(config.Kafka.Address)
	consumer2 := kafkaPackage.NewConsumerKafka()
	orderSyncService2 := order_elastic.NewOrderSyncService(service2, consumer2, producer2, &config, logger)

	logger.Info("Order Elastic Service is starting.")
	go func() {
		if err := orderSyncService1.StartConsumeOrder(); err != nil {
			logger.Fatalf("Order sync service failed, shutting down the server. Error:%v", err)
		}
	}()
	if err := orderSyncService2.StartPushOrder(); err != nil {
		logger.Fatalf("Order sync service failed, shutting down the server. Error:%v", err)
	}
}
