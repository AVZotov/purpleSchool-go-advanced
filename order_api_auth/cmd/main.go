package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"order_api_auth/internal/config"
	"order_api_auth/internal/domain/auth/session"
	"order_api_auth/internal/http/server"
	"order_api_auth/pkg/db"
	pkgLogger "order_api_auth/pkg/logger"
	mw "order_api_auth/pkg/middleware"
	pkgValidator "order_api_auth/pkg/validator"
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

	if err = pkgValidator.Init(); err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Error initializing validator")
		panic(err)
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

	sessionService := session.NewService(nil)

	//TODO: Implement Repository

	session.NewHandler(mux, nil, sessionService)

	stack := mw.Chain(
		mw.RequestIDMiddleware,
		mw.LoggerMiddleware,
	) //TODO: Implement all other mw's

	handler := stack(mux)

	srv := server.New(cfg.HttpServer.Port, handler)
	err = srv.ListenAndServe()
	if err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("failed to start server")
		panic(err)
	}

	//TODO: Delete this

	//TODO: Auth domain setup

	//TODO AUTH Middleware setup
}
