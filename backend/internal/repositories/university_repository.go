package repositories

import (
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/pkg/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertBudget   = "INSERT into universities.budget (prestige, name, rank, speciality, quality, points, dormitory, labs, sport, scholarship, region_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	insertContract = "INSERT into universities.contract (prestige, name, rank, speciality, quality, points, cost, dormitory, labs, sport, scholarship, region_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
)

type UniversityRepository struct {
	pg *pgxpool.Pool
}

func NewUniversityRepository(ctx context.Context, cfg postgres.Config) (*UniversityRepository, error) {
	pool, err := postgres.New(ctx, cfg, "universities")
	if err != nil {
		return nil, fmt.Errorf("NewUniversityRepository: failed to connect to users postgres: %w", err)
	}
	return &UniversityRepository{pg: pool}, nil
}

func (ur *UniversityRepository) FillBudget(ctx context.Context, u *entities.University) error {
	_, err := ur.pg.Exec(ctx, insertBudget,
		u.Prestige, u.Name, u.Rank, u.Speciality,
		u.Quality, u.Points, u.Dormitory,
		u.Labs, u.Sport, u.Scholarship, u.Region)
	if err != nil {
		return fmt.Errorf("FillBudget: %w", err)
	}
	return nil
}

func (ur *UniversityRepository) FillContract(ctx context.Context, u *entities.University) error {
	_, err := ur.pg.Exec(ctx, insertContract,
		u.Prestige, u.Name, u.Rank, u.Speciality,
		u.Quality, u.Points, u.Cost, u.Dormitory,
		u.Labs, u.Sport, u.Scholarship, u.Region)
	if err != nil {
		return fmt.Errorf("FillContract: %w", err)
	}
	return nil
}
