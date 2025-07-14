package main

import (
	"log"
	"order/internal/config"
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
}
