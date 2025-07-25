package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"order/internal/config"
	pkgLogger "order/pkg/logger"
	"time"
)

type DB struct {
	*gorm.DB
}

func New(config *config.Config, appLogger pkgLogger.Logger) (*DB, error) {
	appLogger.Info("Initializing database connection",
		"host", config.Database.Host,
		"port", config.Database.Port,
		"database", config.Database.Name,
	)

	gLogger := NewGormLogger(appLogger, gormLogger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  getGormLogLevel(config.Env.String()),
		IgnoreRecordNotFoundError: true,
		Colorful:                  config.Env.IsDev(),
	})

	dsn := config.Database.PsqlDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gLogger,
	})
	if err != nil {
		appLogger.Error("Failed to connect to database", "error", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	dbStruct := &DB{
		DB: db,
	}

	if err = dbStruct.HealthCheck(appLogger); err != nil {
		return nil, err
	}

	appLogger.Info("Database connection established successfully")

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
		}
	}
	return result.RowsAffected, result.Error
}

func (db *DB) FindAll(models any) error {
	return db.DB.Find(models).Error
}

func (db *DB) UpdatePartial(model any, id uint, fields map[string]any) (int64, error) {
	result := db.DB.Model(model).Where("id = ?", id).Updates(fields)
	return result.RowsAffected, result.Error
}

func (db *DB) UpdateAll(model any) error {
	return db.DB.Save(model).Error
}

func (db *DB) Delete(module any, conditions ...any) (int64, error) {
	result := db.DB.Delete(module, conditions)
	return result.RowsAffected, result.Error
}

func getGormLogLevel(env string) gormLogger.LogLevel {
	switch env {
	case "dev":
		return gormLogger.Warn
	case "prod":
		return gormLogger.Error
	default:
		return gormLogger.Info
	}
}
