package main

import (
	"fmt"
	"order/internal/config"
	"path"
)

const ConfigPath = "./config/env"
const DevFile = "configs.yml"

func main() {
	cfg := config.MustLoadConfig(path.Join(ConfigPath, DevFile))
	fmt.Println(cfg)
}
