package main

import (
	"fmt"
	"order_simple/internal/config"
	"order_simple/pkg/db"
	pkgLogger "order_simple/pkg/logger"
	"order_simple/pkg/migrations"
)

func main() {

	cfg, err := config.MustLoadConfig("configs.yml")
	if err != nil {
		panic(err)
	}

	pkgLogger.Init()
	pkgLogger.LogApplicationStart(cfg)

	database, err := db.New(cfg)
	if err != nil {
		pkgLogger.LogDatabaseError("error to initiate DB", err)
		panic(err)
	}
	err = migrations.RunMigrations(database.DB)
	if err != nil {
		return
	}

	fmt.Println(database.DB.RowsAffected)
}
