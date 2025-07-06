package main

import (
	"fmt"
	"link_shortener/config"
	"link_shortener/internal/http-server/handlers/email/verify"
	"link_shortener/internal/http-server/handlers/system"
	"link_shortener/internal/http-server/router"
	"link_shortener/internal/http-server/server"
	"link_shortener/pkg/security"
	"log"
	"net/http"
)

func main() {
	configs, err := config.NewConfig("mailhog")
	if err != nil {
		log.Fatal(err)
	}
	emailSecrets, err := configs.GetEmailSecrets()
	if err != nil {
		log.Fatal(err)
	}
	router := router.NewRouter()
	err = registerHandlers(router, emailSecrets)
	if err != nil {
		log.Fatal(err)
		return
	}
	server := server.NewServer("8081", router)
	log.Println("Starting server on port 8081")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func registerHandlers(router *http.ServeMux, secrets []byte) error {
	err := verify.NewVerificationHandler(router, secrets, security.NewHash())
	if err != nil {
		return fmt.Errorf("error creating 'NewVerificationHandler' handler: %s", err)
	}
	system.NewHealthCheckHandler(router)
	return nil
}
