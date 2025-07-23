package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var Validator *validator.Validate

func Init() error {
	Validator = validator.New()
	err := Validator.RegisterValidation("session_id", validateSessionID)
	if err != nil {
		return err
	}
	return nil
}

func ValidateStruct(s any) error {
	return Validator.Struct(s)
}

func validateSessionID(fl validator.FieldLevel) bool {
	sessionID := fl.Field().String()

	matched, _ := regexp.MatchString(`^[a-f0-9]{64}`, sessionID)
	return matched
}
