package auth

type RequestForSession struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type RequestForVerification struct {
	SessionID string `json:"sessionId" validate:"required,hexadecimal,len=64"`
	Code      int    `json:"code" validate:"required,gte=1000,lte=9999"`
}

type ResponseWithSessionID struct {
	SessionID string `json:"sessionId" validate:"required,hexadecimal,len=64"`
}

type ResponseWithJWT struct {
	Token string `json:"token" validate:"required"`
}

type Session struct {
	Phone     string `json:"phone,omitempty" validate:"omitempty,e164"`
	SessionID string `json:"sessionId,omitempty" validate:"omitempty,hexadecimal,len=64"`
	SMSCode   int    `json:"code,omitempty" validate:"omitempty,gte=1000,lte=9999"`
}
