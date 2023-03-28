package cmd

import (
	"OrderUserProject/internal/apps/order-elastic"
	"OrderUserProject/internal/apps/order-elastic/handler"
	"OrderUserProject/internal/configs"
	"OrderUserProject/pkg"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func StartOrderElastic() {
	e := echo.New()

	// TODO : echo olmayacak, log normal log.

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

	OrderService := order_elastic.NewOrderElasticService()

	handler.NewOrderElasticHandler(e, *OrderService)

	// TODO : root servis olacak (while true)
	// Start server
	go func() {
		if err := e.Start(config.Server.Port["orderElastic"]); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server!")
		}
	}()

	pkg.GracefulShutdown(e, 10*time.Second)
}
