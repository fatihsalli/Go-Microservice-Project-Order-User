package cmd

import (
	"OrderUserProject/docs"
	"OrderUserProject/internal/apps/order-api"
	"OrderUserProject/internal/apps/order-api/handler"
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/repository"
	"OrderUserProject/pkg"
	kafka_Package "OrderUserProject/pkg/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"time"
)

// @title           Order API
// @version         1.0
// @description     This is a sample restful server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8011
// @BasePath  /api
func StartOrderAPI() {
	e := echo.New()

	// TODO : Kafka producer içeri handlera paslanacak
	// logger instead of echo.log we use 'logrus' package
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	e.Logger = log.Logger()
	e.Use(pkg.Logger())
	log.Info("Logger enabled!!")

	// get config
	config := configs.GetConfig("test")

	// to create kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Errorf("Cannot create a producer: %v", err)
	}
	producer := kafka_Package.NewProducerKafka(p, "orderID-created-v01")

	// to create repo and service
	mongoOrderCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.OrderCollectionName)
	OrderRepository := repository.NewOrderRepository(mongoOrderCollection)
	OrderService := order_api.NewOrderService(OrderRepository)

	// to create handler
	handler.NewOrderHandler(e, OrderService, producer)

	// if we don't use this swagger give an error
	docs.SwaggerInfo.Host = "localhost:8011"
	// add swagger
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	// Start server
	go func() {
		if err := e.Start(config.Server.Port["orderAPI"]); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server!")
		}
	}()

	pkg.GracefulShutdown(e, 10*time.Second)
}
