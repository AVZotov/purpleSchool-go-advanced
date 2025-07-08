package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" required:"true"`
	DbStorage   string `yaml:"db_storage" required:"true"`
	MailService `yaml:"mail_service" required:"true"`
	HTTPServer  `yaml:"http_server" required:"true"`
}

type MailService struct {
	Name     string `yaml:"name" required:"true"`
	Email    string `yaml:"email" required:"true"`
	Password string `yaml:"password" required:"true"`
	Host     string `yaml:"host" required:"true"`
	Port     string `yaml:"port" required:"true"`
	Address  string `yaml:"address" required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"http://localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoadConfig(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
		return nil
	}
	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatal(err)
		return nil
	}

	return &config
}
