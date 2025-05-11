package config

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

// setEnvVars sets env vars for tests
func setEnvVars(t *testing.T, vars map[string]string) func() {
	for key, value := range vars {
		err := os.Setenv(key, value)
		assert.NoError(t, err)
	}
	return func() {
		for key := range vars {
			err := os.Unsetenv(key)
			if err != nil {
				return
			}
		}
	}
}

// createTempEnvFile creates env file for tests
func createTempEnvFile(t *testing.T, content string) (string, func()) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, ".env")
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(t, err)
	return filePath, func() {
		err = os.Remove(filePath)
		if err != nil {
			return
		}
	}
}

// TestNewUserConfig_Success tests success case
func TestNewUserConfig_Success(t *testing.T) {
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

	cfg, err := NewUserConfig()
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

// TestNewUserConfig_WithEnvFile tests env file case
func TestNewUserConfig_WithEnvFile(t *testing.T) {
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

	cfg, err := NewUserConfig()
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

// TestNewUserConfig_EnvFileError tests nonexistent env file case
func TestNewUserConfig_EnvFileError(t *testing.T) {
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

	cfg, err := NewUserConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}
