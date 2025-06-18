package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Token string `env:"Token,required,notEmpty"`
	API   string `env:"API,required,notEmpty"`
}

func Must(cfg *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return cfg
}

func NewFromEnv() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config parsing failed: %w", err)
	}

	return &cfg, nil
}
