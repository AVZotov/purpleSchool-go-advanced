package session

import (
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Phone     string `json:"phone" validate:"required,e164" gorm:"index"`
	SessionID string `json:"session_id,omitempty" validate:"omitempty,hexadecimal,len=64" gorm:"index;size:64"`
	SMSCode   string `json:"sms_code,omitempty" validate:"omitempty,len=4,numeric" gorm:"size:4"`
}
