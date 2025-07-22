package session

import (
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Phone     string `json:"phone" validate:"required,phone" gorm:"index"`
	SessionID string `json:"session_id,omitempty" validate:"required,session_id" gorm:"index;size:64"`
	SMSCode   string `json:"sms_code,omitempty" validate:"required,sms_code" gorm:"size:4"`
}
