package order

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"order_api_cart/pkg/db"
	"order_api_cart/pkg/db/models"
	pkgErr "order_api_cart/pkg/errors"
	pkgLogger "order_api_cart/pkg/logger"
)

type Repository struct {
	DB *db.DB
}

func NewRepository(db *db.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) findOrCreateUser(ctx context.Context) (*models.User, error) {
	var user models.User
	phone, err := parsePhone(ctx)
	if err != nil {
		return nil, err
	}
	result := r.DB.FirstOrCreate(&user, models.User{Phone: phone})
	if result.Error != nil {
		pkgLogger.ErrorWithRequestID(ctx, "Failed to find or create user", logrus.Fields{
			"error": result.Error.Error(),
			"phone": phone,
		})
		return nil, pkgErr.ErrQueryFailed
	}

	if result.RowsAffected > 0 {
		pkgLogger.InfoWithRequestID(ctx, "New user created", logrus.Fields{
			"user_id": user.ID,
			"phone":   phone,
		})
	} else {
		pkgLogger.InfoWithRequestID(ctx, "Existing user found", logrus.Fields{
			"user_id": user.ID,
			"phone":   phone,
		})
	}

	return &user, nil
}

func (r *Repository) getProducts(ctx context.Context, productIDs []uint) ([]models.Product, error) {
	var products []models.Product

	result := r.DB.Where("id IN ?", productIDs).Find(&products)
	if result.Error != nil {
		pkgLogger.ErrorWithRequestID(ctx, "Failed to find products", logrus.Fields{
			"error": result.Error.Error(),
		})
		return nil, pkgErr.ErrQueryFailed
	}

	if len(products) != len(productIDs) {
		pkgLogger.ErrorWithRequestID(ctx, "Failed to find products", logrus.Fields{})
		return nil, pkgErr.ErrProductNotFound
	}

	return products, nil
}

func (r *Repository) createOrder(ctx context.Context, order *models.Order) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			pkgLogger.ErrorWithRequestID(ctx, "Failed to create order", logrus.Fields{
				"error": err.Error(),
			})
			return err
		}

		pkgLogger.InfoWithRequestID(ctx, "Order created in transaction", logrus.Fields{
			"order_id":       order.ID,
			"user_id":        order.UserID,
			"products_count": len(order.Products),
		})

		return nil
	})
	if err != nil {
		pkgLogger.ErrorWithRequestID(ctx, "Failed to create order", logrus.Fields{
			"error": err.Error(),
		})
		return pkgErr.ErrTransactionFailed
	}

	return nil
}
