package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"order/internal/config"
	"order/internal/db_models/product"
)

type DB struct {
	*gorm.DB
}

func New(config *config.Config) (*DB, error) {
	db, err := gorm.Open(postgres.Open(config.Database.PsqlDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) RunMigrations() error {
	return db.AutoMigrate(
		&product.Product{},
	)
}
