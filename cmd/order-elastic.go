package cmd

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func StartOrderElastic() {

	// Logger instead of standard log we use 'logrus' package
	log := logrus.StandardLogger()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	log.Info("Logger enabled!!")

}
