package repositories

import (
	"University-Selection-Service/pkg/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticRepository struct {
	pg *pgxpool.Pool
}

func NewAnalyticRepository(ctx context.Context, cfg postgres.Config) (*AnalyticRepository, error) {
	pool, err := postgres.New(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("NewAnalyticRepository: failed to connect to user postgres: %w", err)
	}
	return &AnalyticRepository{pg: pool}, nil
}
