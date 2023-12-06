package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "../.env"

type (
	Config struct {
		HTTPServer HTTPServer
		RabbitMQ   RabbitMQ
		MongoDB    MongoDB
		SMTP       SMTP
	}
	HTTPServer struct {
		Address     string        `env:"SERVER_ADDRESS"`
		Timeout     time.Duration `env:"SERVER_TIMEOUT"`
		IdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT"`
	}
	RabbitMQ struct {
		Host     string `env:"RABBITMQ_HOST"`
		User     string `env:"RABBITMQ_USER"`
		Password string `env:"RABBITMQ_PASSWORD"`
		Queue    string `env:"RABBITMQ_QUEUE"`
		Port     string `env:"RABBITMQ_PORT"`
	}
	MongoDB struct {
		Host     string `env:"MONGODB_HOST"`
		User     string `env:"MONGODB_USER"`
		Password string `env:"MONGODB_PASSWORD"`
		Port     string `env:"MONGODB_PORT"`
	}
	SMTP struct {
		Login    string `env:"SMTP_LOGIN"`
		Password string `env:"SMTP_PASSWORD"`
		Provider string `env:"SMTP_PROVIDER"`
		Port     string `env:"SMTP_PORT"`
	}
)

func NewConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
