package config

type EmailSecrets struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required,min=1,max=65535"`
	Address  string `json:"address" validate:"required"`
	Provider string `json:"provider" validate:"required"`
}
