package main

import (
	"log"
	"order/internal/config"
	"order/internal/http_server/router"
	"order/internal/http_server/server"
	"order/pkg/db"
	"os"
)

const DevFile = "configs.yml"

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Application panicked: %v", rec)
			os.Exit(1)
		}
	}()

	cfg, err := config.MustLoadConfig(DevFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	database, err := db.New(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = database.RunMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	r := router.New(database)
	s := server.New(cfg.HttpServer.Port, r)

	log.Printf("Server starting on port %s", cfg.HttpServer.Port)
	if err = s.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
