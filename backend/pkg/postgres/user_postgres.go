package postgres

import (
	"University-Selection-Service/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST"`
	Port     string `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT"`
	Username string `yaml:"POSTGRES_USER" env:"POSTGRES_USER"`
	Password string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB"`

	MinConns int32 `yaml:"MIN_CONNS" env:"POSTGRES_MIN_CONN"`
	MaxConns int32 `yaml:"MAX_CONNS" env:"POSTGRES_MAX_CONN"`
}

func New(ctx context.Context, c Config) (*pgxpool.Pool, error) {
	log := logger.GetLoggerFromCtx(ctx)
	connstring := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_min_conns=%d&pool_max_conns=%d",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.MinConns,
		c.MaxConns)
	conn, err := pgxpool.New(ctx, connstring)
	if err != nil {
		return nil, fmt.Errorf("New: failed to connect to postgres: %w", err)
	} else {
		log.Info(ctx, "connected to user_postgres")
	}

	migration, err := migrate.New(
		"file://db/user_db_migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			c.Username,
			c.Password,
			c.Host,
			c.Port,
			c.Database,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("New: failed to create migration instance: %w", err)
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("New: failed to Up migration: %w", err)
	}
	log.Info(ctx, "Successfully Applied Migration")
	return conn, nil
}
