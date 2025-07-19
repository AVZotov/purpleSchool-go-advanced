package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"order_simple/internal/config"
	"time"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	logger.Info("starting application")

	cfg, err := config.MustLoadConfig("configs.yml")
	if err != nil {
		logger.Fatal(err)
	}

}
