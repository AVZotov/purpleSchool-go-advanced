package migrations

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	pkgLogger "order_api_cart/pkg/logger"
)

func RunMigrations(db *gorm.DB) error {
	pkgLogger.Logger.Info("start migrations")

	models := []any{} // TODO: Add models for migrations

	if err := db.AutoMigrate(models...); err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to run migrations")
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
