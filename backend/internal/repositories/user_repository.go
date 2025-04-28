package repositories

import (
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/pkg/postgres"
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
	GetByLoginSQLRequest                   = "SELECT * FROM schema_name.users WHERE login=$1"
	SaveRefreshTokenSQLRequest             = "INSERT INTO schema_name.refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)"
	CreateUserSQLRequest                   = "INSERT INTO schema_name.users (login, password) VALUES ($1, $2, $3) RETURNING Id"
	RevokeAllActiveTokensForUserSQLRequest = "DELETE FROM schema_name.refresh_tokens WHERE user_id = $1"
	GetUserIDByRefreshTokenSQLRequest      = "DELETE FROM schema_name.refresh_tokens WHERE token = $1"
)

func NewUserRepository(ctx context.Context, cfg postgres.Config) (*UserRepository, error) {
	pool, err := postgres.New(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("NewUserRepository: failed to connect to user postgres: %w", err)
	}
	return &UserRepository{pg: pool}, nil
}

func (ur *UserRepository) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	user := &entities.User{}
	queryRow := ur.pg.QueryRow(ctx, GetByLoginSQLRequest, login)
	err := queryRow.Scan(&user.Id, &user.Login, &user.Password, &user.Ege, &user.Gpa,
		&user.Speciality, &user.EduType, &user.Town, &user.Financing)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "User not found")
	} else if err != nil {
		return nil, fmt.Errorf("GetUserByLogin: failed to query user by login: %w", err)
	}
	return user, nil
}

func (ur *UserRepository) SaveRefreshToken(ctx context.Context, id int, token string) error {
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err := ur.pg.Exec(ctx, SaveRefreshTokenSQLRequest, id, token, expiresAt)
	if err != nil {
		return fmt.Errorf("SaveRefreshToken: failed to save refresh token: %w", err)
	}
	return nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *entities.User) (int, error) {
	var id int
	queryRow := ur.pg.QueryRow(ctx, CreateUserSQLRequest, user.Login, user.Password)
	err := queryRow.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: failed to query user: %w", err)
	}
	return id, nil
}

func (ur *UserRepository) RevokeAllActiveTokensForUser(ctx context.Context, userId int) error {
	_, err := ur.pg.Exec(ctx, RevokeAllActiveTokensForUserSQLRequest, userId)
	if err != nil {
		return fmt.Errorf("RevokeAllActiveTokensForUser: failed to revoke active tokens for user: %w", err)
	}
	return nil
}

func (ur *UserRepository) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (int, error) {
	var id int
	queryRow := ur.pg.QueryRow(ctx, GetUserIDByRefreshTokenSQLRequest, refreshToken)
	err := queryRow.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("GetUserIDByRefreshToken: failed to query user by refresh token: %w", err)
	}
	return id, nil
}
