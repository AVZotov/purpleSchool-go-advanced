package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	pkgLogger "order_simple/pkg/logger"
	"strings"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(r *http.Request, value any) error {
	err := v.validate.Struct(value)
	if err != nil {
		pkgLogger.WarnWithRequestID(r, "validation failed", logrus.Fields{
			"error": err.Error(),
			"type":  pkgLogger.ValidationError,
		})
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

func (v *Validator) ValidateImageURLs(r *http.Request, value *any) error {
	if p.Images == nil || len(value.Images) == 0 {
		return nil
	}

	for i, imageURL := range p.Images {
		if strings.TrimSpace(imageURL) == "" {
			continue
		}

		if _, err := url.ParseRequestURI(imageURL); err != nil {
			pkgLogger.WarnWithRequestID(r, "invalid image url", logrus.Fields{
				"type":  pkgLogger.ValidationError,
				"error": err.Error(),
			})
			return fmt.Errorf("invalid image URL at position %d (%s): %w", i, imageURL, err)
		}

		if parsed, _ := url.Parse(imageURL); parsed.Scheme == "" {
			pkgLogger.WarnWithRequestID(r, "missing URL scheme", logrus.Fields{
				"type":  pkgLogger.ValidationError,
				"error": errors.New("missing URL scheme").Error(),
			})
			return fmt.Errorf("image URL at position %d missing scheme (http/https): %s", i, imageURL)
		}
	}

	return nil
}
