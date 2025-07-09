package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type MailService struct {
	name     string `yaml:"name" required:"true"`
	email    string `yaml:"email" required:"true"`
	password string `yaml:"password" required:"true"`
	host     string `yaml:"host" required:"true"`
	port     string `yaml:"port" required:"true"`
	address  string `yaml:"address" required:"true"`
}

func (m MailService) Name() string {
	return m.name
}

func (m MailService) Email() string {
	return m.email
}

func (m MailService) Password() string {
	return m.password
}

func (m MailService) Host() string {
	return m.host
}

func (m MailService) Port() string {
	return m.port
}

func (m MailService) Address() string {
	return m.address
}

type HttpServer struct {
	address     string        `yaml:"address" env-default:"http://localhost:8080"`
	port        string        `yaml:"port" env-default:"8080"`
	timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	idleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func (h HttpServer) Address() string {
	return h.address
}

func (h HttpServer) Port() string {
	return h.port
}

func (h HttpServer) Timeout() time.Duration {
	return h.timeout
}

func (h HttpServer) IdleTimeout() time.Duration {
	return h.idleTimeout
}

type Config struct {
	env         string `yaml:"env" required:"true"`
	dbStorage   string `yaml:"db_storage" required:"true"`
	MailService `yaml:"mail_service" required:"true"`
	HttpServer  `yaml:"http_server" required:"true"`
}

func (c Config) Env() string {
	return c.env
}

func (c Config) DbStorage() string {
	return c.dbStorage
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
