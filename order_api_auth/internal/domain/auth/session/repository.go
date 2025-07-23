package session

type Repository interface {
	CreateSession(session *Session) error
	GetSession(sessionID string) (*Session, error)
	DeleteSession(sessionID string) error
}
