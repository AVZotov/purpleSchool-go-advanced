package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"os"
)

type EmailSecrets struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
	Address  string `validate:"required"`
}

func newEmailConfig(envTags map[string]string) (_ *EmailSecrets, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in 'newEmailConfig': %w", err)
		}
	}()
	err = godotenv.Load()
	if err != nil {
		return nil, err
	}
	config := EmailSecrets{
		Email:    os.Getenv(envTags["email"]),
		Password: os.Getenv(envTags["password"]),
		Address:  os.Getenv(envTags["address"]),
	}
	err = validateStruct(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func validateStruct(config *EmailSecrets) error {
	validate := validator.New()
	return validate.Struct(config)
}
