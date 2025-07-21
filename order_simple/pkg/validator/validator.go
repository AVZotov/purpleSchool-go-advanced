package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(value any) error {
	err := v.validate.Struct(value)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
