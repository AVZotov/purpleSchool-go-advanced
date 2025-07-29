package migrations

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	pkgModels "order_api_cart/pkg/db/models"
	pkgErrors "order_api_cart/pkg/errors"
	pkgLogger "order_api_cart/pkg/logger"
)

func RunMigrations(db *gorm.DB) error {
	pkgLogger.Logger.Info("start migrations")

	models := []any{
		pkgModels.User{},
		pkgModels.Order{},
		pkgModels.Product{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(pkgErrors.ErrMigrationFailed.Error())
		return fmt.Errorf("%w: %v", pkgErrors.ErrMigrationFailed, err)
	}

	return nil
}
