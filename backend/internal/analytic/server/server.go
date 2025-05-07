package analytic

import (
	"University-Selection-Service/internal/config"
	"University-Selection-Service/internal/repositories"
	"University-Selection-Service/pkg/api"
	"context"
	"fmt"
)

type Server struct {
	api.AnalyticServer
	rep *repositories.AnalyticRepository
}

func New(ctx context.Context, cfg *config.AnalyticConfig) (*Server, error) {
	r, err := repositories.NewAnalyticRepository(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("NewAnalyticRepository: %w", err)
	}
	return &Server{
		rep: r,
	}, nil
}
