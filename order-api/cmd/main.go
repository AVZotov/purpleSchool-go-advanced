package main

import (
	"log"
	"order/internal/config"
	"order/internal/http_server/router"
	"order/internal/http_server/server"
	"order/pkg/db"
	"path"
)

const ConfigPath = "./config/env"
const DevFile = "configs.yml"

func main() {
	cfg := config.MustLoadConfig(path.Join(ConfigPath, DevFile))
	_, err := db.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := router.New()
	s := server.New(cfg.HttpServer.Port, r)
	log.Print("Listening on port ", cfg.HttpServer.Port)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
