package db

import (
	"errors"
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

func (db *DB) FindById(model any, id uint) (int64, error) {
	result := db.DB.First(model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		} else {
			return 0, fmt.Errorf("finding model failed: %w", result.Error)
		}
	}

	return result.RowsAffected, nil
}

func (db *DB) FindAll(model any) error {
	return db.DB.Find(model).Error
}

func (db *DB) UpdatePartial(model any, id uint, fields map[string]any) (int64, error) {
	result := db.DB.Model(model).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		} else {
			return 0, fmt.Errorf("updating model failed: %w", result.Error)
		}
	}

	return result.RowsAffected, nil
}

func (db *DB) UpdateAll(model any) error {
	return db.DB.Save(model).Error
}

func (db *DB) Delete(module any, conditions ...any) (int64, error) {
	result := db.DB.Delete(module, conditions)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		} else {
			return 0, fmt.Errorf("deleting model failed: %w", result.Error)
		}
	}

	return result.RowsAffected, nil
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
