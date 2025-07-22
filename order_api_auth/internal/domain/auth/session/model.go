package session

import (
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Phone     string `json:"phone" validate:"required,phone"`
	SessionID string `json:"session_id,omitempty" validate:"required,session_id"`
	SMSCode   string `json:"sms_code,omitempty" validate:"required,sms_code"`
	Used      string `json:"used,omitempty"`
}
