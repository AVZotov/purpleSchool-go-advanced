package session

type SendCodeRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type VerifyCodeRequest struct {
	SessionID string `json:"sessionId" validate:"required,hexadecimal,len=64"`
	Code      int    `json:"code" validate:"required,gte=1000,lte=9999"`
}

type ResponseWithSession struct {
	SessionID string `json:"sessionId" validate:"required,hexadecimal,len=64"`
}

type ResponseWithJWT struct {
	Token string `json:"token" validate:"required"`
}
