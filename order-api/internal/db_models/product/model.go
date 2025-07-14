package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Image       pq.StringArray `json:"image,omitempty"`
}

func New(name string, description string, images ...string) (*Product, error) {
	p := &Product{
		Name:        name,
		Description: description,
		Image:       images,
	}

}
