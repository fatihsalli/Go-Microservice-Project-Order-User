package main

import (
	"OrderUserProject/docs"
	"OrderUserProject/internal/apps/order-api"
	handlerOrder "OrderUserProject/internal/apps/order-api/handler"
	"OrderUserProject/internal/apps/user-api"
	handlerUser "OrderUserProject/internal/apps/user-api/handler"
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/repository"
	"OrderUserProject/pkg"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
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

// @host      localhost:8080
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

	mongoUserCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.UserCollectionName)

	mongoOrderCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.OrderCollectionName)

	OrderRepository := repository.NewOrderRepository(mongoOrderCollection)
	UserRepository := repository.NewUserRepository(mongoUserCollection)

	UserService := user_api.NewUserService(*UserRepository)
	OrderService := order_api.NewOrderService(*OrderRepository)

	// to create new app
	handlerUser.NewUserHandler(e, UserService)
	handlerOrder.NewOrderHandler(e, OrderService)

	// if we don't use this swagger give an error
	docs.SwaggerInfo.Host = "localhost:8080"
	// add swagger
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(config.Server.Port))
}
