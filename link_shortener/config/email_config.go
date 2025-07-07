package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"link_shortener/pkg/validate"
	"os"
	"strconv"
	"strings"
)

const (
	EMAIL    = "EMAIL"
	PASSWORD = "PASSWORD"
	HOST     = "HOST"
	PORT     = "PORT"
	ADDRESS  = "ADDRESS"
)

func providers() map[string]string {
	return map[string]string{
		"yandex":  "YANDEX",
		"google":  "GOOGLE",
		"mailhog": "MAILHOG",
	}
}

func newEmailConfig(provider string) (_ *EmailSecrets, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in 'newEmailConfig': %w", err)
		}
	}()
	var emailSecrets *EmailSecrets
	emailSecrets, err = loadEnvSecrets(provider)
	if err != nil {
		return nil, err
	}

	err = validate.StructValidator(emailSecrets)
	if err != nil {
		return nil, err
	}

	return emailSecrets, nil
}

func getValidProvider(provider string) (string, error) {
	p, ok := providers()[strings.ToLower(provider)]
	if !ok {
		err := errors.New("email provider not supported")
		return "", fmt.Errorf("error in 'getValidProvider': %w", err)
	}
	return p, nil
}

func loadEnvSecrets(provider string) (_ *EmailSecrets, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in 'loadEnvSecrets': %w", err)
		}
	}()
	err = godotenv.Load()
	if err != nil {
		return nil, err
	}

	provider, err = getValidProvider(provider)
	if err != nil {
		return nil, err
	}

	var port int
	port, err = getEnvInt(fmt.Sprintf("%s_%s", provider, PORT))
	if err != nil {
		return nil, err
	}

	secrets := EmailSecrets{
		Email:    os.Getenv(fmt.Sprintf("%s_%s", provider, EMAIL)),
		Password: os.Getenv(fmt.Sprintf("%s_%s", provider, PASSWORD)),
		Host:     os.Getenv(fmt.Sprintf("%s_%s", provider, HOST)),
		Port:     port,
		Address:  os.Getenv(fmt.Sprintf("%s_%s", provider, ADDRESS)),
		Provider: strings.ToLower(provider),
	}
	return &secrets, nil
}

func getEnvInt(key string) (int, error) {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0, fmt.Errorf("error in 'getEnvInt': %w", err)
	}

	return value, err
}
