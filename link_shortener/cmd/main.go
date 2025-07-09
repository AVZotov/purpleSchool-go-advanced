package main

import (
	config "link_shortener/internal/config"
	"link_shortener/internal/http-server/handlers/email/info"
	"link_shortener/internal/http-server/handlers/system"
	router "link_shortener/internal/http-server/router"
	httpserver "link_shortener/internal/http-server/server"
	"link_shortener/pkg/security"
	storage "link_shortener/pkg/storage/local_storage"
	"log"
	"log/slog"
	"net/http"

	l "link_shortener/pkg/logger"

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

	configPath := path.Join(ConfigPath, DevFile)
	configs := config.MustLoadConfig(configPath)

	logger := l.NewLogger(configs.Env())

	localStorage, err := storage.New(configs.Env(), logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	router := router.NewRouter()

	err = registerHandlers(router, configs, localStorage, logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	server := httpserver.NewServer(configs.HttpServer.Port(), router)
	err = server.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info("Starting server on port %s", configs.HttpServer.Port())
}

func registerHandlers(router *http.ServeMux, configs *config.Config,
	storage *storage.Storage, logger *slog.Logger) error {
	const fn = "link_shortener.cmd.main.registerHandlers"
	logger.With(fn)
	emailSecrets := &configs.MailService
	err := email.New(router, emailSecrets, security.NewHashHandler(), storage)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = info.NewInfoHandler(router, secrets)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	system.NewHealthCheckHandler(router)
	return nil
}
