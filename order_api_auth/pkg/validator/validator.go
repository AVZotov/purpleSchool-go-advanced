package validator

import (
	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func Init() error {
	Validator = validator.New()
	return nil
}

func ValidateStruct(s any) error {
	return Validator.Struct(s)
}
