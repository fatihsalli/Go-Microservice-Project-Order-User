package main

import (
	"OrderUserProject/docs"
	"OrderUserProject/internal/apps/order-api"
	handler_order "OrderUserProject/internal/apps/order-api/handler"
	"OrderUserProject/internal/apps/user-api"
	handler_user "OrderUserProject/internal/apps/user-api/handler"
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/repository"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	e.Logger.Fatal(e.Start(config.Server.Port))
}

// TODO: Echo loglaması ezilerek logrus kullanılacak
// TODO: Middleware
// TODO: Graceful shutdown araştırılacak
