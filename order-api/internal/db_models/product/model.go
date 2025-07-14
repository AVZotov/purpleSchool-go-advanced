package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"order/pkg/validator"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Images      pq.StringArray `json:"image,omitempty" validate:"url"`
}

func New(name string, description string, images ...string) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Images:      images,
	}
}

func (p *Product) Validate() error {
	v := validator.New()
	return v.Validate(p)
}
