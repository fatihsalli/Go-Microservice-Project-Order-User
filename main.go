package main

import (
	"OrderUserProject/docs"
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/controller"
	"OrderUserProject/internal/repository"
	"OrderUserProject/internal/service"
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

	mongoUserCollection := configs.ConnectDB(config.Database.Connection).
		Database(config.Database.DatabaseName).Collection(config.Database.UserCollectionName)

	// to create new repository with singleton pattern
	UserRepository := repository.GetSingleInstancesRepository(mongoUserCollection)

	// to create new service with singleton pattern
	UserService := service.GetSingleInstancesService(UserRepository)

	// to create new app
	controller.NewUserHandler(e, UserService)

	// if we don't use this swagger give an error
	docs.SwaggerInfo.Host = "localhost:8080"
	// add swagger
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(config.Server.Port))
}
