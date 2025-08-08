package migrations

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	pkgModels "order_api_cart/pkg/db/models"
	pkgErr "order_api_cart/pkg/errors"
	pkgLogger "order_api_cart/pkg/logger"
)

func RunMigrations(db *gorm.DB) error {
	pkgLogger.Logger.Info("start migrations")

	models := []any{
		pkgModels.User{},
		pkgModels.Order{},
		pkgModels.Product{},
		pkgModels.Session{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(pkgErr.ErrServiceUnavailable.Error())
		return fmt.Errorf("%w: %v", pkgErr.ErrServiceUnavailable, err)
	}

	return nil
}
