package validate

import "github.com/go-playground/validator/v10"

func StructValidator(str any) error {
	validate := validator.New()
	return validate.Struct(str)
}
