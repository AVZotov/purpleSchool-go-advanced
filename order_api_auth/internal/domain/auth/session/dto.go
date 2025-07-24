package session

type ResponseWithSession struct {
	SessionID string `json:"session_id" validate:"required,hexadecimal,len=64"`
}

type ResponseWithJWT struct {
	JWT string `json:"token" validate:"required"`
}
