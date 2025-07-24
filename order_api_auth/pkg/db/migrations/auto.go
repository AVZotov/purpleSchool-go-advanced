package migrations

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"order_api_auth/internal/domain/auth/session"
	pkgLogger "order_api_auth/pkg/logger"
)

func RunMigrations(db *gorm.DB) error {
	pkgLogger.Logger.Info("start migrations")

	models := []any{
		&session.Session{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to run migrations")
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
