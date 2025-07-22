package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"order_api_auth/internal/config"
	pkgLogger "order_api_auth/pkg/logger"
	"os"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Application panicked: %v", rec)
			os.Exit(1)
		}
	}()

	const DevFile = "configs.yml"

	cfg, err := config.MustLoadConfig(DevFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	pkgLogger.Init()
	pkgLogger.Logger.WithFields(logrus.Fields{
		"logger": pkgLogger.Logger.Level,
		"env":    cfg.Env.String(),
	}).Info("logger initialized")

	//TODO: DB setup

	//TODO: Handler setup

	//TODO AUTH Middleware setup
}
