package config

import (
	"University-Selection-Service/pkg/postgres"
)

type Config struct {
	postgres.Config `yaml:"POSTGRES" env:"POSTGRES"`
}
