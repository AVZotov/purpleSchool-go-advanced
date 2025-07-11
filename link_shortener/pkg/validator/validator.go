package validator

import "github.com/go-playground/validator/v10"

type StructValidator struct{}

func (s StructValidator) Validate(str any) error {
	validate := validator.New()
	return validate.Struct(str)
}
