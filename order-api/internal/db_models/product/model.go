package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"net/url"
	"order/pkg/validator"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Images      pq.StringArray `json:"image,omitempty" gorm:"type:text[]"`
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
	if err := v.Validate(p); err != nil {
		return err
	}

	for _, imageURL := range p.Images {
		if imageURL == "" {
			continue
		}
		if _, err := url.Parse(imageURL); err != nil {
			return err
		}
	}
	return nil
}
