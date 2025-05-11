package config

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewAnalyticCfg_Success tests success case
func TestNewAnalyticCfg_Success(t *testing.T) {
	unset := setEnvVars(t, map[string]string{
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_DBNAME":   "dbname",
		"INT_PORT":          "50051",
		"REST_PORT":         "8080",
		"JWT_SECRET":        "secret",
	})
	defer unset()

	cfg, err := NewAnalyticCfg()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, "5432", cfg.Postgres.Port)
	assert.Equal(t, "user", cfg.Postgres.Username)
	assert.Equal(t, "password", cfg.Postgres.Password)
	assert.Equal(t, "50051", cfg.INTPort)
	assert.Equal(t, "8080", cfg.RESTPort)
	assert.Equal(t, "secret", cfg.JWTSecret)
}

// TestNewAnalyticCfg_WithEnvFile tests env file case
func TestNewAnalyticCfg_WithEnvFile(t *testing.T) {
	envContent := `
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DBNAME=dbname
INT_PORT=50051
REST_PORT=8080
JWT_SECRET=secret
`
	envFile, cleanup := createTempEnvFile(t, envContent)
	defer cleanup()

	err := godotenv.Load(envFile)
	assert.NoError(t, err)

	cfg, err := NewAnalyticCfg()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, "5432", cfg.Postgres.Port)
	assert.Equal(t, "user", cfg.Postgres.Username)
	assert.Equal(t, "password", cfg.Postgres.Password)
	assert.Equal(t, "50051", cfg.INTPort)
	assert.Equal(t, "8080", cfg.RESTPort)
	assert.Equal(t, "secret", cfg.JWTSecret)
}

// TestNewAnalyticCfg_EnvFileError tests nonexistent env file case
func TestNewAnalyticCfg_EnvFileError(t *testing.T) {
	err := godotenv.Load("nonexistent.env")
	assert.Error(t, err)

	unset := setEnvVars(t, map[string]string{
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "password",
		"INT_PORT":          "50051",
		"REST_PORT":         "8080",
		"JWT_SECRET":        "secret",
	})
	defer unset()

	cfg, err := NewAnalyticCfg()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}
