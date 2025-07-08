package main

import (
	"fmt"
	cfg "link_shortener/internal/config"
	"link_shortener/internal/http-server/handlers/email/info"
	"link_shortener/internal/http-server/handlers/email/verify"
	"link_shortener/internal/http-server/handlers/system"
	r "link_shortener/internal/http-server/router"
	httpserver "link_shortener/internal/http-server/server"
	"link_shortener/pkg/security"
	st "link_shortener/pkg/storage/local_storage"
	"log"
	"log/slog"
	"net/http"

	l "link_shortener/pkg/logger"

	"path"
)

const ConfigPath = "./config/env"
const DevFile = "dev.yml"

//const ProdFile = "prod.yml"

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Recovered from panic: %v", rec)
		}
	}()

	configPath := path.Join(ConfigPath, DevFile)
	configs := cfg.MustLoadConfig(configPath)

	logger := l.NewLogger(configs.Env)

	storage, err := st.NewStorage(configs.Env, logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	router := r.NewRouter()

	err = registerHandlers(router, configs, storage, logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//TODO: add port to configs
	server := httpserver.NewServer("8081", router)
	err = server.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info("Starting server on port 8081")

}

func registerHandlers(router *http.ServeMux, configs *cfg.Config, storage *st.Storage, logger *slog.Logger) error {
	logger.With("link_shortener.cmd.registerHandlers()")
	emailSecrets := &configs.MailService
	err := verify.NewVerificationHandler(router, emailSecrets, security.NewHashHandler(), storage)
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
