package postgres

import (
	"University-Selection-Service/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`

	MinConns int32 `env:"POSTGRES_MIN_CONNS"`
	MaxConns int32 `env:"POSTGRES_MAX_CONNS"`
}

func New(ctx context.Context, conf Config) (*pgxpool.Pool, error) {
	log := logger.GetLoggerFromCtx(ctx)
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_min_conns=%d&pool_max_conns=%d",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.MinConns,
		conf.MaxConns)

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("new: failed to connect to intetgration postgres: %w", err)
	}
	migration, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			conf.Username,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Database,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("new: failed to create migration instance: %w", err)
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("new: failed to Up migration: %w", err)
	}
	log.Info(ctx, "Successfully Applied Migration")
	return conn, nil
}
