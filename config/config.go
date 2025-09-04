package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Server struct {
	Host string `yaml:"host" validate:"required"`
	Port int    `yaml:"port" validate:"required"`
}

func (s Server) Address() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
}

type Database struct {
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	DbName   string `yaml:"dbname" validate:"required"`
	SslMode  string `yaml:"sslmode" validate:"required"`
	Schema   string `yaml:"schema" validate:"required"`
}

func (d Database) Dsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&search_path=%s",
		d.User, d.Password, d.Host, d.Port, d.DbName, d.SslMode, d.Schema,
	)
}

type Config struct {
	Server          Server   `yaml:"server"`
	Database        Database `yaml:"database"`
	TokenSigningKey string   `yaml:"token_signing_key" validate:"required"`
}

func Load() (*Config, error) {
	f, err := os.Open(defineConfigPath())
	if err != nil {
		return nil, errors.WithMessage(err, "open config file")
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err, "close config file")
		}
	}()

	cfg := new(Config)
	err = yaml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "decode config yml file")
	}
	err = validator.New().Struct(cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "validate config")
	}

	return cfg, nil
}

func defineConfigPath() string {
	isDev := os.Getenv("APP_MODE") == "dev"
	if isDev {
		return "./config/config_dev.yml"
	}
	return "./config.yml"
}
