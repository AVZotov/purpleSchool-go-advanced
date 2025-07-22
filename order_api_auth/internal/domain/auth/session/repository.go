package session

type Repository interface {
	CreateSession(session *Session) error
	GetSession(sessionID string) (*Session, error)
	DeleteUsedSession(sessionID string) error
}
