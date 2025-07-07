package main

import (
	"link_shortener/config"
	"link_shortener/internal/http-server/handlers/email/info"
	"link_shortener/internal/http-server/handlers/email/verify"
	"link_shortener/internal/http-server/handlers/system"
	r "link_shortener/internal/http-server/router"
	s "link_shortener/internal/http-server/server"
	l "link_shortener/pkg/logger"
	"link_shortener/pkg/security"
	"link_shortener/pkg/storage/local_storage"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Recovered from panic: %v", rec)
		}
	}()

	logger := l.NewLogger("dev")

	configs, err := config.NewConfig("mailhog")
	if err != nil {
		logger.Error(err.Error())
	}

	storage, err := local_storage.NewStorage("dev")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	emailSecrets, err := configs.GetEmailSecrets()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	router := r.NewRouter()

	err = registerHandlers(router, emailSecrets, storage, logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	server := s.NewServer("8081", router)
	err = server.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info("Starting server on port 8081")

}

func registerHandlers(router *http.ServeMux, secrets []byte, storage *local_storage.Storage, logger *slog.Logger) error {
	logger.With("link_shortener.cmd.registerHandlers()")
	err := verify.NewVerificationHandler(router, secrets, security.NewHash(), storage)
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
