package config

import (
	"University-Selection-Service/pkg/postgres"
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
)

type AnalyticConfig struct {
	Postgres  postgres.Config `env:"POSTGRES"`
	INTPort   string          `env:"INT_PORT"`
	RESTPort  string          `env:"REST_PORT"`
	JWTSecret string          `env:"JWT_SECRET"`
}

func NewAnalyticCfg() (*AnalyticConfig, error) {

	err := godotenv.Load()
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return nil, fmt.Errorf("NewUserConfig: error: %w", err)
	}

	var cfg AnalyticConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("NewUserConfig: reading env error: %w", err)
	}
	return &cfg, nil
}
