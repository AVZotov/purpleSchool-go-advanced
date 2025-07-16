package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"order/internal/config"
	"order/internal/domain/product"
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

func (db *DB) CreateWithLogging(product *product.Product) error {
	return db.DB.Create(product).Error
}

func getGormLogLevel(env string) gormLogger.LogLevel {
	switch env {
	case "dev":
		return gormLogger.Info
	case "prod":
		return gormLogger.Error
	default:
		return gormLogger.Warn
	}
}
