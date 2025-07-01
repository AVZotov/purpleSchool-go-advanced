package config

type Config struct {
	EmailConfig *EmailConfig
}

func NewConfig() (*Config, error) {
	emailConfig, err := newEmailConfig()
	if err != nil {
		return nil, err
	}
	return &Config{
		EmailConfig: emailConfig,
	}, nil
}

func (c *Config) GetEmailConfig() *map[string]string {
	return &map[string]string{
		"email":    c.EmailConfig.Email,
		"password": c.EmailConfig.Password,
		"address":  c.EmailConfig.Email,
	}
}
