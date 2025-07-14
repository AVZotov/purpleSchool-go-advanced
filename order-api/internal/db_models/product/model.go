package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Image       pq.StringArray `json:"image"`
}

func New(name string, description string, images ...string) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Image:       images,
	}
}
