package session

type Service interface {
	SendCode(phone string) (sessionID string, err error)
	VerifyCode(sessionID, code string) (jwt string, err error)
}
