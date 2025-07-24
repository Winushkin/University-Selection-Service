package repository

import (
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/pkg/postgres"
	"University-Selection-Service/pkg/security"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type UserRepository struct {
	pg *pgxpool.Pool
}

var (
	//go:embed sql/get_user_by_login.sql
	GetByLoginRequest string

	//go:embed sql/save_refresh_token.sql
	SaveRefreshTokenRequest string

	//go:embed sql/create_user.sql
	CreateUserRequest string

	//go:embed sql/revoke_all_active_tokens_for_user.sql
	RevokeAllActiveTokensForUserRequest string

	//go:embed sql/get_user_id_by_refresh_token.sql
	GetUserIDByRefreshTokenRequest string

	//go:embed sql/get_user_by_id.sql
	GetUserByIDRequest string

	//go:embed sql/fill_info.sql
	FillInfoRequest string
)

// NewUserRepository returns new user repository with connection to DB
func NewUserRepository(ctx context.Context, cfg postgres.Config) (*UserRepository, error) {
	pool, err := postgres.New(ctx, cfg, "users")
	if err != nil {
		return nil, fmt.Errorf("NewUserRepository: failed to connect to users postgres: %w", err)
	}
	return &UserRepository{pg: pool}, nil
}

// GetByLogin returns user by his login
func (ur *UserRepository) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	user := &entities.User{}
	queryRow := ur.pg.QueryRow(ctx, GetByLoginRequest, login)
	err := queryRow.Scan(&user.Id, &user.Login, &user.Password, &user.Ege,
		&user.Speciality, &user.Town, &user.Financing)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "User not found")
	} else if err != nil {
		return nil, fmt.Errorf("GetUserByLogin: failed to query users by login: %w", err)
	}
	return user, nil
}

// GetByID returns user by his ID
func (ur *UserRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	user := &entities.User{}
	queryRow := ur.pg.QueryRow(ctx, GetUserByIDRequest, id)
	err := queryRow.Scan(&user.Id, &user.Login, &user.Password, &user.Ege,
		&user.Speciality, &user.Town, &user.Financing)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "User not found")
	} else if err != nil {
		return nil, fmt.Errorf("GetUserByID: failed to query users by id: %w", err)
	}
	return user, nil
}

// SaveRefreshToken saves refresh token in DB
func (ur *UserRepository) SaveRefreshToken(ctx context.Context, id int, token string) error {
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err := ur.pg.Exec(ctx, SaveRefreshTokenRequest, id, token, expiresAt)
	if err != nil {
		return fmt.Errorf("SaveRefreshToken: failed to save refresh token: %w", err)
	}
	return nil
}

// CreateUser inserts user into DB
func (ur *UserRepository) CreateUser(ctx context.Context, user *entities.User) (int, error) {
	var id int
	queryRow := ur.pg.QueryRow(ctx, CreateUserRequest, user.Login, security.HashPassword(user.Password))
	err := queryRow.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: failed to query users: %w", err)
	}
	return id, nil
}

// RevokeAllActiveTokensForUser revokes active users refresh tokens in DB
func (ur *UserRepository) RevokeAllActiveTokensForUser(ctx context.Context, userId int) error {
	_, err := ur.pg.Exec(ctx, RevokeAllActiveTokensForUserRequest, userId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("RevokeAllActiveTokensForUser: failed to revoke active tokens for users: %w", err)
	}
	return nil
}

// GetUserIDByRefreshToken returns user ID by his refresh token
func (ur *UserRepository) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (int, error) {
	var id int
	queryRow := ur.pg.QueryRow(ctx, GetUserIDByRefreshTokenRequest, refreshToken)
	err := queryRow.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("GetUserIDByRefreshToken: failed to query users by refresh token: %w", err)
	}
	return id, nil
}

// FillInfo inserts user information into DB
func (ur *UserRepository) FillInfo(ctx context.Context, user *entities.User) error {
	_, err := ur.pg.Exec(ctx, FillInfoRequest, user.Ege, user.Speciality,
		user.Town, user.Financing, user.Id)
	if err != nil {
		return fmt.Errorf("FillInfo: failed to fill info: %w", err)
	}
	return nil
}
