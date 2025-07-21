package product

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	pkgLogger "order_simple/pkg/logger"
	pkgValidator "order_simple/pkg/validator"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Images      pq.StringArray `json:"images,omitempty" gorm:"type:text[]" validate:"http_url"`
}

func (p *Product) String() string {
	return fmt.Sprintf("%s: %s: %+v ", p.Name, p.Description, p.Images)
}

func (p *Product) TableName() string {
	return "products"
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

func (p *Product) Validate(r *http.Request) error {
	validator := pkgValidator.New()
	err := validator.Validate(p)
	if err != nil {
		pkgLogger.ErrorWithRequestID(r, "error validating product", logrus.Fields{
			"type":  pkgLogger.ValidationError,
			"error": err.Error(),
		})
		return fmt.Errorf("error validating product: %w", err)
	}

	return nil
}
