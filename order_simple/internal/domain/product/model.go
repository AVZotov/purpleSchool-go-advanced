package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	pkgValidator "order_simple/pkg/validator"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Images      pq.StringArray `json:"images,omitempty" gorm:"type:text[]"`
}

func (p *Product) ToFieldsMap() map[string]interface{} {
	fields := make(map[string]interface{})

	if p.Name != "" {
		fields["name"] = p.Name
	}
	if p.Description != "" {
		fields["description"] = p.Description
	}
	if p.Images != nil {
		fields["images"] = p.Images
	}

	return fields
}

func (p *Product) HasFields() bool {
	return p.Name != "" || p.Description != "" || p.Images != nil
}

func (p *Product) BeforeCreate(_ *gorm.DB) error {
	return pkgValidator.Validator.Validate()
}
