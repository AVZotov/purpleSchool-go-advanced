package main

import (
	"log"
	"order/internal/config"
	"order/internal/http_server/router"
	"order/internal/http_server/server"
	"order/pkg/db"
	"path"
)

const DevFile = "configs.yml"

func main() {
	cfg := config.MustLoadConfig(path.Join(DevFile))

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
