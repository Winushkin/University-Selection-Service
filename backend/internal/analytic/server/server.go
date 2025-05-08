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
	rep     *repositories.AnalyticRepository
	userCli api.UserServiceClient
}

func New(ctx context.Context, cfg *config.AnalyticConfig, client api.UserServiceClient) (*Server, error) {
	r, err := repositories.NewAnalyticRepository(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("NewAnalyticRepository: %w", err)
	}
	return &Server{
		rep:     r,
		userCli: client,
	}, nil
}

func (s *Server) Analyze(ctx context.Context, request *api.AnalyzeRequest) (*api.AnalyzeResponse, error) {
	id, ok := ctx.Value("id").(int)
	if !ok {
		return nil, fmt.Errorf("id is not int")
	}
	_, err := s.userCli.ProfileDataForAnalytic(ctx, &api.ProfileDataForAnalyticRequest{
		Id: int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("ProfileDataForAnalytic: %w", err)
	}
	return nil, nil
}
