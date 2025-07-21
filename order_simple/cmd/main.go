package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"order_simple/internal/config"
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
			log.Printf("application panick recovered: %v", rec)
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
	stack := middleware.Chain(
		middleware.RequestIDMiddleware,
		middleware.LoggerMiddleware,
	)
	stack(mux)

	server := server.New(cfg.HttpServer.Port, stack)

	fmt.Println(database.DB.RowsAffected)
}
