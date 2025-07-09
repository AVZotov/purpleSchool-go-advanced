package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type MailService struct {
	Name     string `yaml:"name" env-required:"true"`
	Email    string `yaml:"email" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	Address  string `yaml:"address" env-required:"true"`
}

func (m MailService) GetName() string {
	return m.Name
}

func (m MailService) GetEmail() string {
	return m.Email
}

func (m MailService) GetPassword() string {
	return m.Password
}

func (m MailService) GetHost() string {
	return m.Host
}

func (m MailService) GetPort() string {
	return m.Port
}

func (m MailService) GetAddress() string {
	return m.Address
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"http://localhost:8080"`
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func (h HttpServer) GetAddress() string {
	return h.Address
}

func (h HttpServer) GetPort() string {
	return h.Port
}

func (h HttpServer) GetTimeout() time.Duration {
	return h.Timeout
}

func (h HttpServer) GetIdleTimeout() time.Duration {
	return h.IdleTimeout
}

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	DbStorage   string `yaml:"db_storage" env-required:"true"`
	MailService `yaml:"mail_service" env-required:"true"`
	HttpServer  `yaml:"http_server" env-required:"true"`
}

func (c Config) GetEnv() string {
	return c.Env
}

func (c Config) GetDbStorage() string {
	return c.DbStorage
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
