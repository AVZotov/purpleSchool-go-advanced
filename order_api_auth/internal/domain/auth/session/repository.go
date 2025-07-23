package session

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"order_api_auth/pkg/db"
	pkgLogger "order_api_auth/pkg/logger"
)

type Repository interface {
	CreateSession(*http.Request, *Session) error
	GetSession(sessionID string) (*Session, error)
	DeleteSession(sessionID string) error
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
		pkgLogger.ErrorWithRequestID(r, "error storing session in DB", logrus.Fields{})
		return fmt.Errorf("error storing session in DB: %w", err)
	}

	return nil
}

func (rep *RepositorySession) GetSession(sessionID string) (*Session, error) {
	//TODO: Implement me
	return nil, nil
}

func (rep *RepositorySession) DeleteSession(sessionID string) error {
	//TODO: Implement me
	return nil
}
