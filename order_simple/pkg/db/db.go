package db

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"order_simple/internal/config"
	pkgLogger "order_simple/pkg/logger"
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

func (db *DB) Create(r *http.Request, v any) error {
	err := db.DB.Create(v).Error
	if err != nil {
		pkgLogger.ErrorWithRequestID(r, "creation failed", logrus.Fields{
			"type": pkgLogger.DBError,
			"err":  err.Error(),
		})
		return fmt.Errorf("creation failed: %w", err)
	}

	return nil
}

func (db *DB) FindById(r *http.Request, model any, id uint) (int64, error) {
	result := db.DB.First(model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		} else {
			pkgLogger.ErrorWithRequestID(r, "finding model failed", logrus.Fields{
				"type": pkgLogger.DBError,
				"err":  result.Error.Error(),
			})
			return 0, fmt.Errorf("finding model failed: %w", result.Error)
		}
	}

	return result.RowsAffected, nil
}

func (db *DB) UpdatePartial(r *http.Request, model any, id uint, fields map[string]any) (int64, error) {
	result := db.DB.Model(model).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		} else {
			pkgLogger.ErrorWithRequestID(r, "updating partial failed", logrus.Fields{
				"type": pkgLogger.DBError,
				"err":  result.Error.Error(),
			})
		}
	}
	return result.RowsAffected, nil
}

func (db *DB) Delete(r *http.Request, module any, conditions ...any) (int64, error) {
	result := db.DB.Delete(module, conditions)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		} else {
			pkgLogger.ErrorWithRequestID(r, "deleting model failed", logrus.Fields{
				"type": pkgLogger.DBError,
				"err":  result.Error.Error(),
			})
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
