package migrations

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"order_simple/internal/domain/product"
	pkgLogger "order_simple/pkg/logger"
	"time"
)

func RunMigrations(db *gorm.DB) error {
	start := time.Now()
	pkgLogger.Logger.WithFields(logrus.Fields{
		"type": pkgLogger.DBMigration,
	}).Info("starting migrations")

	models := []any{
		&product.Product{},
	}

	err := db.AutoMigrate(models...)
	if err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"type":  pkgLogger.DBError,
		}).Error("failed to run migrations")
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	pkgLogger.Logger.WithFields(logrus.Fields{
		"type":         pkgLogger.DBMigration,
		"duration_ms":  time.Since(start).Milliseconds(),
		"models_count": len(models),
	}).Info("finished migrations")
	return nil
}
