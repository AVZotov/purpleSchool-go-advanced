package config

type Config struct {
	EmailConfig *EmailSecrets
}

func NewConfig(envTags map[string]string) (*Config, error) {
	emailConfig, err := newEmailConfig(envTags)
	if err != nil {
		return nil, err
	}
	return &Config{
		EmailConfig: emailConfig,
	}, nil
}

func (c *Config) GetGmailSecrets() *map[string]string {
	return &map[string]string{
		"email":    c.EmailConfig.Email,
		"password": c.EmailConfig.Password,
		"address":  c.EmailConfig.Address,
	}
}
