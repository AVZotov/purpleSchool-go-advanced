package main

import (
	"fmt"
	"link_shortener/internal/config"
	"link_shortener/internal/http-server/handlers/email/info"
	"link_shortener/internal/http-server/handlers/email/verify"
	"link_shortener/internal/http-server/handlers/system"
	t "link_shortener/internal/http-server/handlers/types"
	r "link_shortener/internal/http-server/router"
	httpserver "link_shortener/internal/http-server/server"
	"link_shortener/pkg/logger"
	"link_shortener/pkg/security"
	storage "link_shortener/pkg/storage/local_storage"
	storageLoggerWrapper "link_shortener/pkg/storage/logger"
	"link_shortener/pkg/validator"
	"log"
	"net/http"
	"path"
)

const ConfigPath = "./config/env"
const DevFile = "dev.yml"

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Recovered from panic: %v", rec)
		}
	}()
	const fn = "cmd.main.main"

	configPath := path.Join(ConfigPath, DevFile)
	configs := config.MustLoadConfig(configPath)

	slogLogger := logger.NewLogger(configs.Env())
	slogLogger.With(fn)

	handlersLogger := logger.NewWrapper(slogLogger)
	storageLogger := storageLoggerWrapper.New(slogLogger)

	localStorage, err := storage.New(configs.Env(), storageLogger)
	if err != nil {
		slogLogger.Error(fmt.Sprintf("%s: %v", fn, err))
		return
	}

	router := r.NewRouter()

	err = registerHandlers(router, configs, localStorage, handlersLogger)
	if err != nil {
		slogLogger.Error(fmt.Sprintf("%s: %v", fn, err))
		return
	}

	slogLogger.Info("Starting server on port " + configs.HttpServer.Port())

	server := httpserver.NewServer(configs.HttpServer.Port(), router)

	err = server.ListenAndServe()
	if err != nil {
		slogLogger.Error(fmt.Sprintf("%s: %v", fn, err))
		return
	}
}

func registerHandlers(router *http.ServeMux, configs *config.Config,
	localStorage *storage.Storage, logger t.Logger) error {

	const fn = "cmd.main.registerHandlers"
	logger.With(fn)

	structValidator := &validator.StructValidator{}

	hashHandler := security.NewHashHandler()

	verify.New(router, configs.MailService, hashHandler, localStorage, structValidator, logger)

	err := info.New(router, configs.MailService, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %v", fn, err))
		return fmt.Errorf("%s: %w", fn, err)
	}

	system.NewHealthCheckHandler(router)

	logger.Debug("All handlers registered successfully")
	return nil
}
