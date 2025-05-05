package config

import (
	"University-Selection-Service/integrations/postgres"
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
)

type IntegrationConfig struct {
	Postgres    postgres.Config `env:"POSTGRES"`
	BudgetURL   string          `env:"BUDGET_URL"`
	ContractURL string          `env:"CONTRACT_URL"`
}

func New() (*IntegrationConfig, error) {
	err := godotenv.Load("../.env")
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return nil, fmt.Errorf("NewIntegrationConfig: error: %w", err)
	}

	var cfg IntegrationConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("NewIntegrationConfig: reading env error: %w", err)
	}
	return &cfg, nil
}
