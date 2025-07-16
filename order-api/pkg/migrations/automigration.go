package migrations

import (
	"fmt"
	"gorm.io/gorm"
	pkgLogger "order/pkg/logger"
)

func RunMigrations(db *gorm.DB, appLogger pkgLogger.Logger, modules ...any) error {
	appLogger.Debug("starting database migrations")

	err := db.AutoMigrate(modules...)
	if err != nil {
		appLogger.Error("Failed to run migrations", "error", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	appLogger.Info("Database migrations completed successfully")
	return nil
}
