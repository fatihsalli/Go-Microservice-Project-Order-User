package cmd

import (
	"OrderUserProject/internal/apps/user-api"
	"OrderUserProject/internal/apps/user-api/handler"
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/repository"
	"OrderUserProject/pkg"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func StartUserAPI() {
	e := echo.New()

	// Logger instead of echo.log we use 'logrus' package
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	e.Logger = log.Logger()
	e.Use(pkg.Logger())
	log.Info("Logger enabled!!")

	config := configs.GetConfig("test")

	// To create repo and service
	mongoUserCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.UserCollectionName)
	UserRepository := repository.NewUserRepository(mongoUserCollection)
	UserService := user_api.NewUserService(UserRepository)

	// To create new app
	handler.NewUserHandler(e, UserService)

	// Start server
	go func() {
		if err := e.Start(config.Server.Port["userAPI"]); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server!")
		}
	}()

	pkg.GracefulShutdown(e, 10*time.Second)
}
