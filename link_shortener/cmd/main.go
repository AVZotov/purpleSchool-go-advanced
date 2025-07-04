package main

import (
	"fmt"
	"link_shortener/config"
	"link_shortener/internal/handlers/system"
	"link_shortener/internal/handlers/verify"
	"link_shortener/pkg"
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
	router := pkg.NewRouter()
	err = registerHandlers(router, emailSecrets)
	if err != nil {
		log.Fatal(err)
		return
	}
	server := pkg.NewServer("8081", router)
	log.Println("Starting server on port 8081")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func registerHandlers(router *http.ServeMux, secrets []byte) error {
	err := verify.NewEmailHandler(router, secrets, security.NewHash())
	if err != nil {
		return fmt.Errorf("error creating 'NewEmailHandler' handler: %s", err)
	}
	system.NewHealthCheckHandler(router)
	return nil
}
