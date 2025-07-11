package main

import (
	"errors"
	"link_shortener/internal/config"
	"link_shortener/internal/http-server/handlers/email/verify"
	"link_shortener/internal/http-server/handlers/system"
	"link_shortener/internal/http-server/router"
	"link_shortener/internal/http-server/server"
	"link_shortener/pkg/container"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const ConfigPath = "./config/env"
const DevFile = "dev.yml"

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Application panicked: %v", rec)
			os.Exit(1)
		}
	}()

	// Load configuration
	configPath := getConfigPath()
	cfg := config.MustLoadConfig(configPath)

	// Initialize container
	ctr, err := container.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	mux := router.NewRouter()
	err = registerHandlers(mux, ctr)
	if err != nil {
		ctr.Logger.Error("Failed to register handlers: %v", err)
		return
	}

	srv := server.New(cfg.HttpServer.Port, mux)

	ctr.Logger.Info("Starting server",
		"port", cfg.HttpServer.Port,
		"env", cfg.Env)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		ctr.Logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func getConfigPath() string {
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		return configPath
	}
	return filepath.Join(ConfigPath, DevFile)
}

func registerHandlers(mux *http.ServeMux, ctr *container.Container) error {
	err := verify.New(mux, ctr.Logger, ctr.EmailService, ctr.HashService, ctr.Storage, ctr.Validator)
	if err != nil {
		ctr.Logger.Error("Failed to register verification handler:", "error", err)
		return err
	}

	//err := info.New(router, configs.MailService, logger)
	//if err != nil {
	//	logger.Error(fmt.Sprintf("%s: %v", fn, err))
	//	return fmt.Errorf("%s: %w", fn, err)
	//}

	system.NewHealthCheckHandler(mux)

	ctr.Logger.Debug("All handlers registered successfully")
	return nil
}
