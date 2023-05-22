package cmd

import (
	docs "OrderUserProject/docs/user"
	"OrderUserProject/internal/apps/user-api"
	"OrderUserProject/internal/apps/user-api/handler"
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/repository"
	"OrderUserProject/pkg"
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

// @title           User Microservice
// @version         1.0
// @description     This is a user microservice project.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8012
// @BasePath  /api
func StartUserAPI() {
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

	// Get config
	config := configs.GetConfig("test")

	// Create new repo and new service
	mongoUserCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.UserCollectionName)
	UserRepository := repository.NewUserRepository(mongoUserCollection)
	UserService := user_api.NewUserService(UserRepository)

	// Create new app
	handler.NewUserHandler(e, UserService, v)

	// If we don't use this swagger give an error
	docs.SwaggerInfouserAPI.Host = "localhost:8012"
	// Add swagger (InstanceName is important!)
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("userAPI")))

	// Start server as asynchronous
	go func() {
		if err := e.Start(config.Server.Port["userAPI"]); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server!")
		}
	}()

	// Graceful Shutdown
	pkg.GracefulShutdown(e, 10*time.Second)
}
