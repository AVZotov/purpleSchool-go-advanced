package main

import (
	"order/internal/config"
	"path"
)

const ConfigPath = "./config/env"
const DevFile = "configs.yml"

func main() {
	_ = config.MustLoadConfig(path.Join(ConfigPath, DevFile))
}
