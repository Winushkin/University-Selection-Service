package user

import (
	"University-Selection-Service/internal/config"
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/internal/repositories"
	"University-Selection-Service/pkg/api"
	"University-Selection-Service/pkg/security"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
	}

	_, err = s.rep.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "error creating users")
	}

	user, err = s.rep.GetByLogin(ctx, request.Login)

	if err != nil {
		return nil, status.Error(codes.Internal, "error getting users")
	}

	accessToken, err := s.generateAccessToken(user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "error generating access token")
	}
	expAt := time.Now().Add(15 * time.Minute).Unix()

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, status.Error(codes.Internal, "error generating refresh token")
	}

	err = s.rep.SaveRefreshToken(ctx, user.Id, refreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "error saving refresh token")
	}

	return &api.SignUpResponse{
		Access:    accessToken,
		Refresh:   refreshToken,
		ExpiresAt: expAt,
	}, nil
}

func (s *Server) Login(ctx context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {
	user, err := s.rep.GetByLogin(ctx, request.Login)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	if user.Login != request.Login || security.CheckPasswordHash(request.Password, user.Password) {
		return nil, status.Error(codes.FailedPrecondition, "User login and password do not match")
	}

	err = s.rep.RevokeAllActiveTokensForUser(ctx, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "error revoking active tokens for users")
	}

	accessToken, err := s.generateAccessToken(user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "error generating access token")
	}
	expAt := time.Now().Add(15 * time.Minute).Unix()

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, status.Error(codes.Internal, "error generating refresh token")
	}

	err = s.rep.SaveRefreshToken(ctx, user.Id, refreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error saving refresh token: %v", err))
	}

	return &api.LoginResponse{
		Access:    accessToken,
		Refresh:   refreshToken,
		ExpiresAt: expAt,
	}, nil
}

func (s *Server) Refresh(ctx context.Context, request *api.RefreshRequest) (*api.RefreshResponse, error) {
	id, err := s.rep.GetUserIDByRefreshToken(ctx, request.Refresh)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	accessToken, err := s.generateAccessToken(id)
	if err != nil {
		return nil, status.Error(codes.Internal, "error generating access token")
	}
	expAt := time.Now().Add(15 * time.Minute).Unix()

	return &api.RefreshResponse{
		Access:    accessToken,
		ExpiresAt: expAt,
	}, nil
}

func (s *Server) Fill(ctx context.Context, request *api.FillRequest) (*api.FillResponse, error) {
	id := ctx.Value("user_id").(int)
	usr := &entities.User{
		Id:         id,
		Ege:        int(request.Ege),
		Speciality: request.Speciality,
		Town:       request.Town,
		Financing:  request.Financing,
	}
	err := s.rep.FillInfo(ctx, usr)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error filling users: %s", err))
	}
	return &api.FillResponse{}, nil
}

func (s *Server) Logout(ctx context.Context, request *api.LogoutRequest) (*api.LogoutResponse, error) {
	id, err := s.rep.GetUserIDByRefreshToken(ctx, request.Refresh)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	err = s.rep.RevokeAllActiveTokensForUser(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "error revoking active tokens for users")
	}
	return &api.LogoutResponse{}, nil
}

func (s *Server) Profile(ctx context.Context, _ *emptypb.Empty) (*api.ProfileResponse, error) {
	id := ctx.Value("user_id").(int)
	usr, err := s.rep.GetByID(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	return &api.ProfileResponse{
		Ege:        int32(usr.Ege),
		Speciality: usr.Speciality,
		Town:       usr.Town,
		Financing:  usr.Financing,
	}, nil
}

func (s *Server) ProfileDataForAnalytic(ctx context.Context, request *api.ProfileDataForAnalyticRequest) (*api.ProfileDataForAnalyticResponse, error) {
	usr, err := s.rep.GetByID(ctx, int(request.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	return &api.ProfileDataForAnalyticResponse{
		Ege:        int32(usr.Ege),
		Speciality: usr.Speciality,
		Town:       usr.Town,
		Financing:  usr.Financing,
	}, nil
}
