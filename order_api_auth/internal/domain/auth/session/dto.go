package session

type ResponseWithSession struct {
	SessionID string `json:"session_id" validate:"required,hexadecimal,len=64"`
	JWT       string `json:"json,omitempty" validate:"omitempty,jwt"`
}
