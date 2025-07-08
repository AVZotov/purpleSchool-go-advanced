package config

import "time"

type Config struct {
	Env         string `yaml:"env" required:"true"`
	DbStorage   string `yaml:"db_storage" required:"true"`
	MailService `yaml:"mail_service" required:"true"`
	Server      `yaml:"http_server" required:"true"`
}

type MailService struct {
	Email    string `yaml:"email" required:"true"`
	Password string `yaml:"password" required:"true"`
	Host     string `yaml:"host" required:"true"`
	Port     string `yaml:"port" required:"true"`
	Address  string `yaml:"address" required:"true"`
}

type Server struct {
	Address     string        `yaml:"address" set-default:"http://localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" set-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" set-default:"60s"`
}
