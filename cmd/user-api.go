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
	"os"
	"time"
)

func StartUserAPI() {
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

	UserRepository := repository.NewUserRepository(mongoUserCollection)

	UserService := user_api.NewUserService(*UserRepository)

	// to create new app
	handler.NewUserHandler(e, UserService)

	e.Logger.Fatal(e.Start(":8012"))
}
