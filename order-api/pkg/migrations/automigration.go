package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"order/internal/domain/product"
	pkgLogger "order/pkg/logger"
)

func RunMigrations(db *gorm.DB, appLogger pkgLogger.Logger) error {
	appLogger.Debug("starting database migrations")

	models := []any{
		&product.Product{},
	}

	err := db.AutoMigrate(models...)
	if err != nil {
		appLogger.Error("Failed to run migrations", "error", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	appLogger.Info("Database migrations completed successfully")
	return nil
}
