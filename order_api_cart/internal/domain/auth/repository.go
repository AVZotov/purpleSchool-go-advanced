package auth

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"order_api_cart/pkg/db"
	pkgErrors "order_api_cart/pkg/errors"
	pkgLogger "order_api_cart/pkg/logger"
)

type Repository struct {
	DB db.DB
}

func (repo *Repository) FindByOrCreate(ctx context.Context, model any, conditions ...any) error {
	err := repo.DB.FindByOrCrate(model, conditions...)
	if err != nil {
		pkgLogger.ErrorWithRequestID(ctx, pkgErrors.ErrRecordNotFound.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return fmt.Errorf("%w %v", pkgErrors.ErrRecordNotFound, err.Error())
	}

	return nil
}
