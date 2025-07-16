package main

import (
	"log"
	"order/internal/config"
	"order/pkg/container"
	"os"
	"os/signal"
	"syscall"
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

	go func() {
		sygChan := make(chan os.Signal, 1)
		signal.Notify(sygChan, syscall.SIGINT, syscall.SIGTERM)
		<-sygChan

		if err = ctr.Close(); err != nil {
			ctr.Logger.Error("Error during shutdown", "error", err)
		}
		os.Exit(0)
	}()

	ctr.Logger.Info("Starting server")
	err = ctr.Start()
	if err != nil {
		ctr.Logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
