package main

import (
	"order/internal/config"
	"order/pkg/db"
	"path"
)

const ConfigPath = "./config/env"
const DevFile = "configs.yml"

func main() {
	cfg := config.MustLoadConfig(path.Join(ConfigPath, DevFile))
	_, _ = db.New(cfg)
}
