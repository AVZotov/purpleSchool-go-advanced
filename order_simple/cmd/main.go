package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"order_simple/internal/config"
	"order_simple/internal/domain/product"
	"order_simple/internal/http/server"
	"order_simple/pkg/db"
	pkgLogger "order_simple/pkg/logger"
	"order_simple/pkg/middleware"
	"order_simple/pkg/migrations"
	"os"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("application panic recovered: %v", rec)
			os.Exit(1)
		}
	}()

	cfg, err := config.MustLoadConfig("configs.yml")
	if err != nil {
		panic(err)
	}

	pkgLogger.Init()
	pkgLogger.Logger.WithFields(logrus.Fields{
		"env":         cfg.Env.String(),
		"server_host": cfg.HttpServer.Host,
		"server_port": cfg.HttpServer.Port,
		"db_host":     cfg.Database.Host,
		"db_port":     cfg.Database.Port,
	}).Info("configuration loaded")

	database, err := db.New(cfg)
	if err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to connect to database")
		panic(err)
	}

	err = migrations.RunMigrations(database.DB)
	if err != nil {
		return
	}

	mux := http.NewServeMux()

	productRepo := product.NewRepository(database)
	product.New(mux, productRepo)

	stack := middleware.Chain(
		middleware.RequestIDMiddleware,
		middleware.LoggerMiddleware,
	)

	handler := stack(mux)

	srv := server.New(cfg.HttpServer.Port, handler)

	pkgLogger.Logger.WithFields(logrus.Fields{
		"port": cfg.HttpServer.Port,
		"host": cfg.HttpServer.Host,
	}).Info("starting HTTP server")

	if err = srv.ListenAndServe(); err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("failed to start server")
	}
}
