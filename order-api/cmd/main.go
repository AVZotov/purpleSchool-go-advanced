package main

import (
	"log"
	"order/internal/config"
	"order/internal/http_server/router"
	"order/internal/http_server/server"
	"order/pkg/container"
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

	ctr, err := container.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

}
