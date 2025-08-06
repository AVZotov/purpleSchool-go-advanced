package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"order_api_cart/internal/config"
	"order_api_cart/internal/domain/auth"
	"order_api_cart/internal/http/server"
	"order_api_cart/pkg/db"
	"order_api_cart/pkg/db/migrations"
	pkgLogger "order_api_cart/pkg/logger"
	mw "order_api_cart/pkg/middleware"
	pkgValidator "order_api_cart/pkg/validator"
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

	if err = pkgValidator.Init(); err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Error initializing validator")
		panic(err)
	}

	dtb, err := db.New(cfg.Database.PsqlDSN())
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

	if err = migrations.RunMigrations(dtb.DB); err != nil {
		panic(err)
	}
	pkgLogger.Logger.Info("database migration initialized")

	mux := http.NewServeMux()

	authRepo := auth.NewRepository(dtb)
	authService := auth.NewService(authRepo, cfg.JWT.Secret)
	auth.NewHandler(mux, authService)

	//TODO: CART HANDLER WITH MUX WITH JWT MW

	stack := mw.Chain(
		mw.RequestIDMiddleware,
		mw.LoggerMiddleware,
	)

	handler := stack(mux)

	//TODO: EVENTBUS TO DB ADDING USER

	srv := server.New(cfg.HttpServer.Port, handler)
	_ = srv.ListenAndServe()
}
