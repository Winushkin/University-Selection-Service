package user

import (
	"University-Selection-Service/internal/config"
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/internal/repositories"
	"University-Selection-Service/pkg/api"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Server struct {
	api.UserServiceServer
	rep       *repositories.UserRepository
	jwtSecret string
}

func New(ctx context.Context, cfg *config.UserConfig, jwt string) (*Server, error) {
	r, err := repositories.NewUserRepository(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("NewUserRepository: %w", err)
	}
	return &Server{
		rep:       r,
		jwtSecret: jwt,
	}, nil
}

func (s *Server) generateAccessToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *Server) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *Server) SignUp(ctx context.Context, request *api.SignUpRequest) (*api.SignUpResponse, error) {
	user, err := s.rep.GetByLogin(ctx, request.Login)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, fmt.Errorf("SignUp: %w", err)
	}
	if user != nil {
		return nil, status.Error(codes.AlreadyExists, "User already exists")
	}
	user = &entities.User{
		Login:    request.Login,
		Password: request.Password,
		Name:     request.Name,
	}

	_, err = s.rep.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "error creating user")
	}
	return &api.SignUpResponse{}, nil
}

func (s *Server) Login(ctx context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {
	return nil, nil
}

func (s *Server) Refresh(ctx context.Context, request *api.RefreshRequest) (*api.RefreshResponse, error) {
	return nil, nil
}

func (s *Server) Fill(ctx context.Context, request *api.FillRequest) (*api.FillResponse, error) {
	return nil, nil
}

func (s *Server) Edit(ctx context.Context, request *api.EditRequest) (*api.EditResponse, error) {
	return nil, nil
}
