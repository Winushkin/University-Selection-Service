package repositories

import (
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/pkg/postgres"
	"University-Selection-Service/pkg/security"
	"context"
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

const (
	GetByLoginSQLRequest                   = "SELECT id, login, password, coalesce(ege, 0), coalesce(speciality, ''), coalesce(region, ''), coalesce(financing, '')  FROM users.users WHERE login=$1"
	SaveRefreshTokenSQLRequest             = "INSERT INTO users.refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)"
	CreateUserSQLRequest                   = "INSERT INTO users.users (login, password) VALUES ($1, $2) RETURNING id"
	RevokeAllActiveTokensForUserSQLRequest = "DELETE FROM users.refresh_tokens WHERE user_id = $1"
	GetUserIDByRefreshTokenSQLRequest      = "SELECT user_id FROM users.refresh_tokens WHERE token = $1"
	GetByIDSQLRequest                      = "SELECT id, login, password, coalesce(ege, 0), coalesce(speciality, ''), coalesce(region, ''), coalesce(financing, '') FROM users.users WHERE id = $1"
	FillInfoSQLRequest                     = "UPDATE users.users SET ege = $1, speciality = $2, region = $3, financing = $4 WHERE id = $5"
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
	queryRow := ur.pg.QueryRow(ctx, GetByLoginSQLRequest, login)
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
	queryRow := ur.pg.QueryRow(ctx, GetByIDSQLRequest, id)
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
	_, err := ur.pg.Exec(ctx, SaveRefreshTokenSQLRequest, id, token, expiresAt)
	if err != nil {
		return fmt.Errorf("SaveRefreshToken: failed to save refresh token: %w", err)
	}
	return nil
}

// CreateUser inserts user into DB
func (ur *UserRepository) CreateUser(ctx context.Context, user *entities.User) (int, error) {
	var id int
	queryRow := ur.pg.QueryRow(ctx, CreateUserSQLRequest, user.Login, security.HashPassword(user.Password))
	err := queryRow.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: failed to query users: %w", err)
	}
	return id, nil
}

// RevokeAllActiveTokensForUser revokes active users refresh tokens in DB
func (ur *UserRepository) RevokeAllActiveTokensForUser(ctx context.Context, userId int) error {
	_, err := ur.pg.Exec(ctx, RevokeAllActiveTokensForUserSQLRequest, userId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("RevokeAllActiveTokensForUser: failed to revoke active tokens for users: %w", err)
	}
	return nil
}

// GetUserIDByRefreshToken returns user ID by his refresh token
func (ur *UserRepository) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (int, error) {
	var id int
	queryRow := ur.pg.QueryRow(ctx, GetUserIDByRefreshTokenSQLRequest, refreshToken)
	err := queryRow.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("GetUserIDByRefreshToken: failed to query users by refresh token: %w", err)
	}
	return id, nil
}

// FillInfo inserts user information into DB
func (ur *UserRepository) FillInfo(ctx context.Context, user *entities.User) error {
	_, err := ur.pg.Exec(ctx, FillInfoSQLRequest, user.Ege, user.Speciality,
		user.Town, user.Financing, user.Id)
	if err != nil {
		return fmt.Errorf("FillInfo: failed to fill info: %w", err)
	}
	return nil
}
