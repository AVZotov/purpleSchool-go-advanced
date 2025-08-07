package auth

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"order_api_cart/pkg/db"
	"order_api_cart/pkg/db/models"
	pkgErr "order_api_cart/pkg/errors"
	pkgLogger "order_api_cart/pkg/logger"
)

type RepositoryAuth struct {
	DB *db.DB
}

func NewRepository(DB *db.DB) *RepositoryAuth {
	return &RepositoryAuth{DB: DB}
}

func (repo *RepositoryAuth) CreateSession(ctx context.Context, session *models.Session) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		var existing models.Session
		err := tx.Where("session_id = ?", session.SessionID).First(&existing).Error
		if err == nil {
			pkgLogger.ErrorWithRequestID(ctx, "can't create session, record already exists", logrus.Fields{
				"error": err.Error(),
			})

			return pkgErr.ErrRecordExists
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			pkgLogger.ErrorWithRequestID(ctx, pkgErr.ErrQueryFailed.Error(), logrus.Fields{
				"error": err.Error(),
			})
			return pkgErr.ErrQueryFailed
		}

		if err = tx.Create(session).Error; err != nil {
			return pkgErr.ErrTransactionFailed
		}

		return nil
	})
}

func (repo *RepositoryAuth) DeleteSession(ctx context.Context, session *models.Session) error {
	rowsAffected, err := repo.DB.DeleteBy(&models.Session{}, "session_id = ?", session.SessionID)
	if err != nil {
		pkgLogger.ErrorWithRequestID(ctx, pkgErr.ErrQueryFailed.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return pkgErr.ErrQueryFailed
	}
	if rowsAffected == 0 {
		pkgLogger.ErrorWithRequestID(ctx, pkgErr.ErrRecordNotFound.Error(), logrus.Fields{
			"error": pkgErr.ErrSessionNotFound.Error(),
		})
		return pkgErr.ErrRecordNotFound
	}

	return nil
}

func (repo *RepositoryAuth) GetSession(ctx context.Context, session *models.Session, sessionID string) error {

	if err := repo.DB.FindBy(session, "session_id = ?", sessionID); err != nil {
		pkgLogger.ErrorWithRequestID(ctx, pkgErr.ErrSessionNotFound.Error(), logrus.Fields{
			"error": err.Error(),
		})

		return err
	}

	return nil
}
