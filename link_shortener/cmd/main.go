package main

import (
	"fmt"
	"link_shortener/config"
	"link_shortener/internal/http-server/handlers/email/info"
	"link_shortener/internal/http-server/handlers/email/verify"
	"link_shortener/internal/http-server/handlers/system"
	r "link_shortener/internal/http-server/router"
	s "link_shortener/internal/http-server/server"
	"link_shortener/pkg/security"
	"link_shortener/pkg/storage/local_storage"
	"log"
	"net/http"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Recovered from panic: %v", rec)
		}
	}()

	configs, err := config.NewConfig("mailhog")
	if err != nil {
		log.Fatal(err)
	}
	storage, err := local_storage.NewStorage("dev")
	if err != nil {
		log.Fatal(err)
	}
	emailSecrets, err := configs.GetEmailSecrets()
	if err != nil {
		log.Fatal(err)
	}
	router := r.NewRouter()
	err = registerHandlers(router, emailSecrets, storage)
	if err != nil {
		log.Fatal(err)
		return
	}
	server := s.NewServer("8081", router)
	log.Println("Starting server on port 8081")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func registerHandlers(router *http.ServeMux, secrets []byte, storage *local_storage.Storage) error {
	err := verify.NewVerificationHandler(router, secrets, security.NewHash(), storage)
	if err != nil {
		return fmt.Errorf("error creating 'NewVerificationHandler' handler: %s", err)
	}
	err = info.NewInfoHandler(router, secrets)
	if err != nil {
		return fmt.Errorf("error creating 'NewInfoHandler' handler: %s", err)
	}
	system.NewHealthCheckHandler(router)
	return nil
}
