package session

import (
	"context"
	"github.com/sirupsen/logrus"
	"order_api_auth/pkg/db"
	pkgLogger "order_api_auth/pkg/logger"
)

type Repository interface {
	CreateSession(context.Context, *Session) error
	GetSession(context.Context, *Session, string) error
	DeleteSession(context.Context, *Session) error
}

type RepositorySession struct {
	db *db.DB
}

func NewRepository(db *db.DB) *RepositorySession {
	return &RepositorySession{
		db: db,
	}
}

func (rep *RepositorySession) CreateSession(ctx context.Context, session *Session) error {
	err := rep.db.Create(session)
	if err != nil {
		pkgLogger.ErrorWithRequestID(ctx, ErrCreatingSession.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return err
	}

	return nil
}

func (rep *RepositorySession) GetSession(ctx context.Context, session *Session, sessionID string) error {

	if err := rep.db.FindBy(session, "session_id = ?", sessionID); err != nil {
		pkgLogger.ErrorWithRequestID(ctx, ErrGettingSession.Error(), logrus.Fields{
			"error": err.Error(),
		})

		return err
	}

	return nil
}

func (rep *RepositorySession) DeleteSession(ctx context.Context, session *Session) error {
	rowsAffected, err := rep.db.DeleteBy(&Session{}, "session_id = ?", session.SessionID)
	if err != nil {
		pkgLogger.ErrorWithRequestID(ctx, ErrDeletingSession.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return err
	}
	if rowsAffected == 0 {
		pkgLogger.ErrorWithRequestID(ctx, ErrDeletingSession.Error(), logrus.Fields{
			"error": ErrSessionNotFound.Error(),
		})
		return ErrSessionNotFound
	}

	return nil
}
