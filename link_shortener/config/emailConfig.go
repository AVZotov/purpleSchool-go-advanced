package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type EmailConfig struct {
	Email    string
	Password string
	Address  string
}

func newEmailConfig() (_ *EmailConfig, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in 'newEmailConfig': %w", err)
		}
	}()
	err = godotenv.Load()
	if err != nil {
		return nil, err
	}
	config := EmailConfig{Email: os.Getenv("E_EMAIL"), Password: os.Getenv("E_PASSWORD")}
	err = validateConfig(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func validateConfig(config *EmailConfig) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in 'validateConfig': %w", err)
		}
	}()
	var errs []error
	if config.Email == "" {
		errs = append(errs, errors.New("'email' parameter can't be empty"))
	}
	if config.Password == "" {
		errs = append(errs, errors.New("'password' parameter can't be empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
