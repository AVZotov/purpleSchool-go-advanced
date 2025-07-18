package product

import (
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"net/url"
	"order/pkg/validator"
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

func (p *Product) ValidateImageURLs() error {
	if p.Images == nil || len(p.Images) == 0 {
		return nil
	}

	for i, imageURL := range p.Images {
		if strings.TrimSpace(imageURL) == "" {
			continue
		}

		if _, err := url.ParseRequestURI(imageURL); err != nil {
			return fmt.Errorf("invalid image URL at position %d (%s): %w", i, imageURL, err)
		}

		if parsed, _ := url.Parse(imageURL); parsed.Scheme == "" {
			return fmt.Errorf("image URL at position %d missing scheme (http/https): %s", i, imageURL)
		}
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

func (p *Product) BeforeUpdate(_ *gorm.DB) error {
	return p.Validate()
}
