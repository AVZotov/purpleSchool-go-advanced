package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"order_api_auth/internal/config"
	"order_api_auth/pkg/db"
	pkgLogger "order_api_auth/pkg/logger"
	mw "order_api_auth/pkg/middleware"
	"os"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Application panicked: %v", rec)
			os.Exit(1)
		}
	}()

	const DevCfgFile = "configs.yml"

	cfg, err := config.MustLoadConfig(DevCfgFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	pkgLogger.Init()
	pkgLogger.Logger.WithFields(logrus.Fields{
		"logger": pkgLogger.Logger.Level,
		"env":    cfg.Env.String(),
	}).Info("logger initialized")

	dtb, err := db.New(cfg)
	if err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("failed to init db")
		panic(err)
	}
	pkgLogger.Logger.WithFields(logrus.Fields{
		"db_host": cfg.Database.Host,
		"db_port": cfg.Database.Port,
		"dialect": dtb.Dialector.Name(),
	}).Info("database initialized")

	//TODO: Implement DB migrations

	mux := http.NewServeMux()

	stack := mw.Chain(
		mw.RequestIDMiddleware,
		mw.LoggerMiddleware,
	)

	handler := stack(mux)
	print(handler) //TODO: Delete this

	//TODO: Auth domain setup

	//TODO AUTH Middleware setup
}
