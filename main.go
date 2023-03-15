package main

import (
	"OrderUserProject/cmd"
	"OrderUserProject/docs"
	"OrderUserProject/internal/apps/aggregator/handler"
	"OrderUserProject/internal/configs"
	"OrderUserProject/pkg"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"time"
)

// @title           Echo Monolithic Microservice Project
// @version         1.0
// @description     This is a sample restful server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8010
// @BasePath  /api
func main() {
	e := echo.New()

	// Logger
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	e.Logger = log.Logger()
	e.Use(pkg.Logger())
	log.Info("Logger enabled!!")

	config := configs.GetConfig("test")

	// to create new app
	handler.NewGatewayHandler(e)

	/*	mongoUserCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.UserCollectionName)

		mongoOrderCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.OrderCollectionName)

		OrderRepository := repository.NewOrderRepository(mongoOrderCollection)
		UserRepository := repository.NewUserRepository(mongoUserCollection)

		UserService := user_api.NewUserService(*UserRepository)
		OrderService := order_api.NewOrderService(*OrderRepository)

		// to create new app
		handlerUser.NewUserHandler(e, UserService)
		handlerOrder.NewOrderHandler(e, OrderService)*/

	// if we don't use this swagger give an error
	docs.SwaggerInfo.Host = "localhost:8010"
	// add swagger
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	go cmd.StartOrderAPI()
	go cmd.StartUserAPI()

	// Start server
	go func() {
		if err := e.Start(config.Server.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server!")
		}
	}()

	pkg.GracefulShutdown(e, 10*time.Second)
}
