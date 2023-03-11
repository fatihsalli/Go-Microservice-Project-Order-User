package main

import (
	"OrderUserProject/docs"
	"OrderUserProject/internal/apps/order-api"
	handler_order "OrderUserProject/internal/apps/order-api/handler"
	"OrderUserProject/internal/apps/user-api"
	handler_user "OrderUserProject/internal/apps/user-api/handler"
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/kafka"
	"OrderUserProject/internal/repository"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"os"
	"os/signal"
)

// @title           Echo Restful API
// @version         1.0
// @description     This is a sample restful server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api
func main() {
	e := echo.New()

	config := configs.GetConfig("test")

	mongoUserCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.UserCollectionName)

	mongoOrderCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.OrderCollectionName)

	OrderRepository := repository.NewOrderRepository(mongoOrderCollection)
	UserRepository := repository.NewUserRepository(mongoUserCollection)

	UserService := user_api.NewUserService(*UserRepository)
	OrderService := order_api.NewOrderService(*OrderRepository)

	// to create new app
	handler_user.NewUserHandler(e, UserService)
	handler_order.NewOrderHandler(e, OrderService)

	// if we don't use this swagger give an error
	docs.SwaggerInfo.Host = "localhost:8080"
	// add swagger
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	producer, err := kafka.ConnectProducer()

	if err != nil {
		log.Print("Hata meydana geldi")
	}

	consumer, err := kafka.ConnectConsumer()

	if err != nil {
		log.Print("Hata meydana geldi")
	}

	// SIGINT (Ctrl+C) yakalamak için kanal oluşturma
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Producer example
	go func() {
		for {
			message := &sarama.ProducerMessage{Topic: "test-topic", Value: sarama.StringEncoder("test message-11.03.2023")}
			partition, offset, err := producer.SendMessage(message)
			if err != nil {
				log.Fatalf("Error sending message to Kafka: %s", err.Error())
			}
			fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
		}
	}()

	// Consumer example
	go func() {
		partitionConsumer, err := consumer.ConsumePartition("test-topic", 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error creating partition consumer: %s", err.Error())
		}
		defer partitionConsumer.Close()

		for {
			select {
			case msg := <-partitionConsumer.Messages():
				fmt.Printf("Received message from partition %d at offset %d: %s\n", msg.Partition, msg.Offset, string(msg.Value))
			case <-signals:
				return
			}
		}
	}()

	e.Logger.Fatal(e.Start(config.Server.Port))
}
