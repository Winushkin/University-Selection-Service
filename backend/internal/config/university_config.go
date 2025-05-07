package config

import (
	"University-Selection-Service/pkg/postgres"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type UniversityConfig struct {
	Postgres    postgres.Config `env:"POSTGRES"`
	BudgetURL   string          `env:"BUDGET_URL"`
	ContractURL string          `env:"CONTRACT_URL"`
}

func NewUniversityConfig() (*UniversityConfig, error) {
	err := godotenv.Load("./env/universities.env")
	if err != nil {
		return nil, fmt.Errorf("NewUniversityConfig: error: %w", err)
	}

	var cfg UniversityConfig

	if err = cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("NewUniversityConfig: reading env error: %w", err)
	}
	return &cfg, nil
}
