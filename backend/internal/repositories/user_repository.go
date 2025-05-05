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
	GetByLoginSQLRequest                   = "SELECT Id, Login, Password, coalesce(Ege, 0), coalesce(Gpa, 0), coalesce(Speciality, ''), coalesce(EduType, ''), coalesce(Town, ''), coalesce(Financing, '')  FROM schema_name.users WHERE Login=$1"
	SaveRefreshTokenSQLRequest             = "INSERT INTO schema_name.refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)"
	CreateUserSQLRequest                   = "INSERT INTO schema_name.users (Login, Password) VALUES ($1, $2) RETURNING Id"
	RevokeAllActiveTokensForUserSQLRequest = "DELETE FROM schema_name.refresh_tokens WHERE user_id = $1"
	GetUserIDByRefreshTokenSQLRequest      = "SELECT user_id FROM schema_name.refresh_tokens WHERE token = $1"
	GetByIDSQLRequest                      = "SELECT Id, Login, Password, coalesce(Ege, 0), coalesce(Gpa, 0), coalesce(Speciality, ''), coalesce(EduType, ''), coalesce(Town, ''), coalesce(Financing, '') FROM schema_name.users WHERE Id = $1"
	FillInfoSQLRequest                     = "UPDATE schema_name.users SET Ege = $1, Gpa = $2, Speciality = $3, EduType = $4, Town = $5, Financing = $6 WHERE Id = $7"
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

func (ur *UserRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	user := &entities.User{}
	queryRow := ur.pg.QueryRow(ctx, GetByIDSQLRequest, id)
	err := queryRow.Scan(&user.Id, &user.Login, &user.Password, &user.Ege, &user.Gpa,
		&user.Speciality, &user.EduType, &user.Town, &user.Financing)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "User not found")
	} else if err != nil {
		return nil, fmt.Errorf("GetUserByID: failed to query user by id: %w", err)
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
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
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

func (ur *UserRepository) FillInfo(ctx context.Context, user *entities.User) error {
	_, err := ur.pg.Exec(ctx, FillInfoSQLRequest, user.Ege, user.Gpa, user.Speciality, user.EduType,
		user.Town, user.Financing, user.Id)
	if err != nil {
		return fmt.Errorf("FillInfo: failed to fill info: %w", err)
	}
	return nil
}
