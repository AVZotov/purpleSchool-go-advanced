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
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
)

func (e Environment) String() string {
	return string(e)
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

type Database struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Name     string `yaml:"name" env:"DB_NAME" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE" env-default:"require"`
	DbVolume string `yaml:"db_volume" env:"DB_VOLUME" env-default:"./data"`
}

// PsqlDSN return DSN string for connection to Psql database with credentials
// received from configs
func (d Database) PsqlDSN() string {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
	return dsn
}

type HttpServer struct {
	Schema      string        `yaml:"schema" env:"HTTP_SCHEMA" env-required:"true"`
	Host        string        `yaml:"host" env:"HTTP_HOST" env-required:"true"`
	Port        string        `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	Address     string        `yaml:"address" env:"HTTP_ADDRESS"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type JWT struct {
	Secret string `yaml:"secret" env:"JWT_SECRET" env-required:"true"`
}

type Config struct {
	Env        Environment `yaml:"env" env:"APP_ENV" env-required:"true"`
	Database   Database    `yaml:"data_base" env:"APP_DATABASE" env-required:"true"`
	HttpServer HttpServer  `yaml:"http_server" env:"HTTP_SERVER" env-required:"true"`
	JWT        JWT         `yaml:"jwt" env:"JWT_JWT" env-required:"true"`
}

// MustLoadConfig returns pointer on [Config] with all the credentials required
// to establish connection to database and server
func MustLoadConfig(configPath string) (*Config, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return config, nil
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
