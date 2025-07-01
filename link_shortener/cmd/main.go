package main

import (
	"link_shortener/config"
	"link_shortener/pkg"
	"log"
	"net/http"
)

func main() {
	configs, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	router := pkg.NewRouter()
	registerHandlers(router, configs)
	server := pkg.NewServer("8081", router)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func registerHandlers(router *http.ServeMux, config *config.Config) {
	pkg.NewEmailHandler(router, config)
}
