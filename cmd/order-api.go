package cmd

import (
	docs "OrderUserProject/docs/order"
	"OrderUserProject/internal/apps/order-api"
	"OrderUserProject/internal/apps/order-api/handler"
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/repository"
	"OrderUserProject/pkg"
	"OrderUserProject/pkg/kafka"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"time"
)

// @title           Order Microservice
// @version         1.0
// @description     This is an order microservice project.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:30011
// @BasePath  /api
func StartOrderAPI() {
	// Echo instance
	e := echo.New()

	// Validator instance
	v := validator.New()

	// Logger instead of echo.log we use 'logrus' package
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	e.Logger = log.Logger()
	e.Use(pkg.Logger())
	log.Info("Logger enabled!!")

	// Environment value
	env := os.Getenv("environment")

	// Get config
	config := configs.GetConfig(env)

	// Create Kafka producer
	producer := kafka.NewProducerKafka(config.Kafka.Address)

	// Create repo and service
	mongoOrderCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.OrderCollectionName)
	OrderRepository := repository.NewOrderRepository(mongoOrderCollection)
	OrderService := order_api.NewOrderService(OrderRepository)
	ElasticService := order_api.NewElasticService(&config)

	// Pointer address check
	fmt.Printf("%p\n", OrderService)

	// Create handler
	handler.NewOrderHandler(e, OrderService, producer, &config, v, ElasticService)

	// If we don't use this swagger give an error
	docs.SwaggerInfoorderAPI.Host = "localhost:30011"
	// Add swagger
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("orderAPI")))

	// Start server as asynchronous
	go func() {
		if err := e.Start(config.Server.Port["orderAPI"]); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server!")
		}
	}()

	// Graceful Shutdown
	pkg.GracefulShutdown(e, 10*time.Second)
}
