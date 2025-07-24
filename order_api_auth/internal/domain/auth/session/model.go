package session

import (
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Phone     string `json:"phone,omitempty" validate:"omitempty,e164" gorm:"index"`
	SessionID string `json:"sessionId,omitempty" validate:"omitempty,hexadecimal,len=64" gorm:"uniqueIndex;size:64"`
	SMSCode   int    `json:"code,omitempty" validate:"omitempty,min=1000,max=9999" gorm:""`
}
