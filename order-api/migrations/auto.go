package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"order/internal/config"
	"order/internal/db_models/product"
	"os"
	"path"
)

const ConfigPath = "./config/env"
const DevFile = "configs.yml"

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Application panicked: %v", rec)
			os.Exit(1)
		}
	}()

	DSN := config.MustLoadConfig(path.Join(ConfigPath, DevFile)).Database.PsqlDSN()

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&product.Product{})
	if err != nil {
		panic(err)
	}

}
