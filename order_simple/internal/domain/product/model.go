package product

import (
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"net/url"
	"order_simple/pkg/validator"
	"strings"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Images      pq.StringArray `json:"images,omitempty" gorm:"type:text[]"`
}

func (p *Product) Validate() error {
	v := validator.New()
	if err := v.Validate(p); err != nil {
		return err
	}
	return nil
}

func (r *ReplaceRequest) Validate() error {
	v := validator.New()
	if err := v.Validate(r); err != nil {
		return err
	}
	return nil
}

func (r *ReplaceRequest) ToProduct(id uint) *Product {
	return &Product{
		Model:       gorm.Model{ID: id},
		Name:        r.Name,
		Description: r.Description,
		Images:      r.Images,
	}
}

func (u *UpdateRequest) ToFieldsMap() map[string]interface{} {
	fields := make(map[string]interface{})

	if u.Name != nil {
		fields["name"] = *u.Name
	}
	if u.Description != nil {
		fields["description"] = *u.Description
	}
	if u.Images != nil {
		fields["images"] = *u.Images
	}

	return fields
}

func (u *UpdateRequest) HasFields() bool {
	return u.Name != nil || u.Description != nil || u.Images != nil
}

func (p *Product) BeforeCreate(_ *gorm.DB) error {
	return p.Validate()
}
