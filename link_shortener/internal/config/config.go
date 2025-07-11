package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Environment string

const (
	EnvLoc  = "loc"
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
)

func (e Environment) String() string {
	return string(e)
}

func (e Environment) IsLoc() bool {
	return e == EnvLoc
}

func (e Environment) IsDev() bool {
	return e == EnvDev
}

func (e Environment) IsProd() bool {
	return e == EnvProd
}

func (e Environment) IsTest() bool {
	return e == EnvTest
}

type MailService struct {
	Name     string `yaml:"name" env:"MAIL_NAME" env-required:"true"`
	Email    string `yaml:"email" env:"MAIL_EMAIL" env-required:"true"`
	Password string `yaml:"password" env:"MAIL_PASSWORD"`
	Schema   string `yaml:"schema" env:"MAIL_SCHEMA" env-required:"true"`
	Host     string `yaml:"host" env:"MAIL_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"MAIL_PORT"`
	Address  string `yaml:"address" env:"MAIL_ADDRESS"`
}

type HttpServer struct {
	Schema      string        `yaml:"schema" env:"HTTP_SCHEMA" env-required:"true"`
	Host        string        `yaml:"host" env:"HTTP_HOST" env-required:"true"`
	Port        string        `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	Address     string        `yaml:"address" env:"HTTP_ADDRESS"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type Database struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Name     string `yaml:"name" env:"DB_NAME" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE" env-default:"disable"`
}

func (d Database) PsqlDSN() string {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
	return dsn
}

type Config struct {
	Env         Environment `yaml:"env" env:"APP_ENV" env-required:"true"`
	MailService MailService `yaml:"mail_service"`
	HttpServer  HttpServer  `yaml:"http_server"`
	//Database    Database    `yaml:"database"`
}

func MustLoadConfig(configPath string) *Config {
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return config
}

func loadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return loadFromEnv()
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file does not exist: %s, trying environment variables", configPath)
		return loadFromEnv()
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &config, nil
}

func loadFromEnv() (*Config, error) {
	var config Config
	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}

	return &config, nil
}
