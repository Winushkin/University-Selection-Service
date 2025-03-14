package config

import (
	"University-Selection-Service/pkg/postgres"
)

// no env vars
// add env vars

type Config struct {
	postgres.Config `yaml:"POSTGRES" env:"POSTGRES"`
}
