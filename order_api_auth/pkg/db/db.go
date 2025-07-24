package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"order_api_auth/internal/config"
	pkgLogger "order_api_auth/pkg/logger"
)

type DB struct {
	*gorm.DB
}

func New(config *config.Config) (*DB, error) {
	pkgLogger.Logger.Info("connecting to database")

	dsn := config.Database.PsqlDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"type": pkgLogger.DBError,
			"err":  err.Error(),
		}).Error("failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	dbStruct := &DB{
		DB: db,
	}

	pkgLogger.Logger.Info("successfully connected to database")
	if err = dbStruct.healthCheck(); err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"type": pkgLogger.DBError,
			"err":  err.Error(),
		}).Error("failed healthcheck")
		return nil, fmt.Errorf("failed healthcheck: %w", err)
	}

	return dbStruct, nil
}

func (db *DB) Create(v any) error {
	return db.DB.Create(v).Error
}

func (db *DB) FindBy(model any, conditions ...any) error {
	result := db.DB.Where(conditions[0], conditions[1:]...).First(model)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) DeleteBy(model any, conditions ...any) error {
	result := db.DB.Delete(model, conditions)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) healthCheck() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	if err = sqlDB.Ping(); err != nil {
		return err
	}

	return nil
}
