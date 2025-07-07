package config

import "encoding/json"

type Config struct {
	EmailConfig *EmailSecrets
}

func NewConfig(emailProvider string) (*Config, error) {
	emailConfig, err := newEmailConfig(emailProvider)
	if err != nil {
		return nil, err
	}
	return &Config{
		EmailConfig: emailConfig,
	}, nil
}

func (c *Config) GetEmailSecrets() ([]byte, error) {
	return json.Marshal(c.EmailConfig)
}
