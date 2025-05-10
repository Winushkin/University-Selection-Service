package analytic

import (
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/internal/repositories"
	"University-Selection-Service/pkg/api"
	"context"
	"fmt"
)

type Server struct {
	api.AnalyticServer
	userCli      api.UserServiceClient
	RepInterface repositories.AnalyticRepositoryInterface
}

func New(client api.UserServiceClient, r interface {
	GetUniversitiesBySpeciality(ctx context.Context, specialityName string) ([]*entities.University, error)
}) (*Server, error) {
	return &Server{
		userCli:      client,
		RepInterface: r,
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
