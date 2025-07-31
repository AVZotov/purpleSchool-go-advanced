package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Phone  string `gorm:"unique;not null;index"`
	Orders []Order
}
type Order struct {
	gorm.Model
	UserID   uint      `gorm:"index;not null"`
	User     User      `gorm:"foreignKey:UserID"`
	Products []Product `gorm:"many2many:order_products;"`
}

type Product struct {
	gorm.Model
	Name   string  `gorm:"not null;index"`
	Orders []Order `gorm:"many2many:order_products;"`
}

type Session struct {
	gorm.Model
	Phone     string `gorm:"index"`
	SessionID string `gorm:"uniqueIndex;size:64"`
	SMSCode   int    `gorm:""`
}
