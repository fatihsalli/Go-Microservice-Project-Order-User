package cmd

import (
	"OrderUserProject/internal/apps/order-api"
	"OrderUserProject/internal/apps/order-api/handler"
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

	config := configs.GetConfig("test")

	mongoOrderCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.OrderCollectionName)

	OrderRepository := repository.NewOrderRepository(mongoOrderCollection)

	OrderService := order_api.NewOrderService(*OrderRepository)

	handler.NewOrderHandler(e, OrderService)

	// Start server
	go func() {
		if err := e.Start(config.Server.Port["orderAPI"]); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server!")
		}
	}()

	pkg.GracefulShutdown(e, 10*time.Second)
}
