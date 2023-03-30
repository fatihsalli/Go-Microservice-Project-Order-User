package cmd

import (
	"OrderUserProject/docs"
	"OrderUserProject/internal/apps/order-api"
	"OrderUserProject/internal/apps/order-api/handler"
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/repository"
	"OrderUserProject/pkg"
	"OrderUserProject/pkg/kafka"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"time"
)

func StartOrderAPI() {
	e := echo.New()

	// Logger instead of echo.log we use 'logrus' package
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	e.Logger = log.Logger()
	e.Use(pkg.Logger())
	log.Info("Logger enabled!!")

	// Get config
	config := configs.GetConfig("test")

	// Create Kafka producer
	producer := kafka.NewProducerKafka(config.Kafka.Address)

	// Create repo and service
	mongoOrderCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.OrderCollectionName)
	OrderRepository := repository.NewOrderRepository(mongoOrderCollection)
	OrderService := order_api.NewOrderService(OrderRepository)

	// Create handler
	handler.NewOrderHandler(e, OrderService, producer, &config)

	// If we don't use this swagger give an error
	docs.SwaggerInfo.Host = "localhost:8011"
	// Add swagger
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	// Start server
	go func() {
		if err := e.Start(config.Server.Port["orderAPI"]); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server!")
		}
	}()

	// Graceful Shutdown
	pkg.GracefulShutdown(e, 10*time.Second)
}
