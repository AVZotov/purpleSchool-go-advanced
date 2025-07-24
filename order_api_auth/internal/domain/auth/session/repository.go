package session

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"order_api_auth/pkg/db"
	pkgLogger "order_api_auth/pkg/logger"
)

type Repository interface {
	CreateSession(*http.Request, *Session) error
	GetSession(*http.Request, *Session, string) error
	DeleteSession(*http.Request, *Session) error
}

type RepositorySession struct {
	db *db.DB
}

func NewRepository(db *db.DB) *RepositorySession {
	return &RepositorySession{
		db: db,
	}
}

func (rep *RepositorySession) CreateSession(r *http.Request, session *Session) error {
	err := rep.db.Create(session)
	if err != nil {
		pkgLogger.ErrorWithRequestID(r, ErrCreatingSession.Error(), logrus.Fields{
			"error": err.Error(),
		})
		return err
	}

	return nil
}

func (rep *RepositorySession) GetSession(r *http.Request, session *Session, sessionID string) error {

	if err := rep.db.FindBy(session, "session_id = ?", sessionID); err != nil {
		pkgLogger.ErrorWithRequestID(r, ErrGettingSession.Error(), logrus.Fields{
			"error": err.Error(),
		})

		return err
	}

	return nil
}

func (rep *RepositorySession) DeleteSession(r *http.Request, session *Session) error {
	if err := rep.db.DeleteBy(
		&Session{}, "session_id = ?", session.SessionID); err != nil {
		pkgLogger.ErrorWithRequestID(r, ErrDeletingSession.Error(), logrus.Fields{
			"error": err.Error(),
		})

		return err
	}

	return nil
}
