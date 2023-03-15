package pkg

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(instance *echo.Echo, timeout time.Duration) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGINT)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Infof("Shutting down server with {%v}", timeout)

	if err := instance.Shutdown(ctx); err != nil {
		log.Errorf("Error while shutting down: %v", err)
	} else {
		log.Info("Server was shut down gracefully")
	}
}
