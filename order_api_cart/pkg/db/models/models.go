package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Phone  string  `gorm:"unique;not null;index" json:"phone"`
	Orders []Order `json:"orders"`
}
type Order struct {
	gorm.Model
	UserID   uint      `gorm:"index;not null" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID" json:"user"`
	Products []Product `gorm:"many2many:order_products;" json:"products"`
}

type Product struct {
	gorm.Model
	Name   string  `gorm:"not null;index" json:"name"`
	Orders []Order `gorm:"many2many:order_products;" json:"orders"`
}
