package product

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

func New(name string, description string) *Product {
	return &Product{
		Name:        name,
		Description: description,
	}
}
