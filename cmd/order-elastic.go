package cmd

import (
	order_elastic "OrderUserProject/internal/apps/order-elastic"
	"OrderUserProject/internal/configs"
	kafkaPackage "OrderUserProject/pkg/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func StartOrderElastic() {

	// Logger instead of standard log we use 'logrus' package
	logger := logrus.StandardLogger()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logger.Info("Logger enabled!!")

	// Get config
	config := configs.GetConfig("test")

	orderElasticService := order_elastic.NewOrderElasticService(&config)

	// To create kafka producer as a 'ProducerKafka' struct
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		logger.Errorf("Cannot create a producer: %v", err)
	}
	producer := kafkaPackage.NewProducerKafka(p, config.Elasticsearch.TopicName["OrderModel"])

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Errorf("Kafka consumer didn't work. Error:%v", err)
	}
	consumer := kafkaPackage.NewConsumerKafka(c)
	orderElasticSync := order_elastic.NewOrderSyncService(orderElasticService, consumer, producer, &config)

	logger.Info("Order Elastic Service is starting.")
	if err := orderElasticSync.StartConsumeOrder(); err != nil {
		logger.Fatalf("Order sync service failed, shutting down the server. Error:%v", err)
	}
	go func() {
		if err := orderElasticSync.StartPushOrder(); err != nil {
			logger.Fatalf("Order sync service failed, shutting down the server. Error:%v", err)
		}
	}()
}
