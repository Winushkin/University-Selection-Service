package repositories

import (
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/pkg/logger"
	"University-Selection-Service/pkg/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"strings"
)

const (
	insertSpecialityRequest      = "Insert INTO universities.specialities (university_id, name, budget_points, contract_points, cost) VALUES ($1, $2, $3, $4, $5)"
	insertUniversityRequest      = "INSERT INTO universities.universities (name, site, prestige, rank, quality, scholarship, dormitory, labs, sport, region_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	getRegionIdByNameRequest     = "SELECT id FROM universities.regions WHERE name = $1"
	getUniversityIdByNameRequest = "SELECT id FROM universities.universities WHERE name = $1"
)

type UniversityRepoInterface interface {
	FillRegions(ctx context.Context, regions []string) error
	InsertUniversity(ctx context.Context, u *entities.University) error
	InsertSpeciality(ctx context.Context, s *entities.Speciality) error
}

type UniversityRepository struct {
	pg *pgxpool.Pool
}

// NewUniversityRepository returns new university repository with connection to DB
func NewUniversityRepository(ctx context.Context, cfg postgres.Config) (*UniversityRepository, error) {
	pool, err := postgres.New(ctx, cfg, "universities")
	if err != nil {
		return nil, fmt.Errorf("NewUniversityRepository: failed to connect to users postgres: %w", err)
	}
	return &UniversityRepository{pg: pool}, nil
}

// GetRegionIdByName returns region ID by his name
func (ur *UniversityRepository) GetRegionIdByName(ctx context.Context, name string) (int, error) {
	var id int
	queryRow := ur.pg.QueryRow(ctx, getRegionIdByNameRequest, name)
	err := queryRow.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("GetRegionIdByName: %s: %w", name, err)
	}
	return id, nil
}

// GetUniversityIdByName returns university ID by his name
func (ur *UniversityRepository) GetUniversityIdByName(ctx context.Context, name string) (int, error) {
	var id int
	queryRow := ur.pg.QueryRow(ctx, getUniversityIdByNameRequest, name)
	err := queryRow.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("GetUniversityIdByName: %s: %w", name, err)
	}
	return id, nil
}

// FillRegions fills regions into DB
func (ur *UniversityRepository) FillRegions(ctx context.Context, regions []string) error {
	var valuesString []string
	var valuesArgs []interface{}
	for i, region := range regions {
		valuesString = append(valuesString, fmt.Sprintf("($%d)", i+1))
		valuesArgs = append(valuesArgs, region)
	}

	insertRegionsRequest :=
		fmt.Sprintf("INSERT INTO universities.regions (name) VALUES %s ON CONFLICT (name) DO NOTHING ",
			strings.Join(valuesString, ","))
	_, err := ur.pg.Exec(ctx, insertRegionsRequest, valuesArgs...)

	if err != nil {
		return fmt.Errorf("FillRegions: %w", err)
	}
	return nil
}

// InsertUniversity insert universities into DB
func (ur *UniversityRepository) InsertUniversity(ctx context.Context, u *entities.University) error {
	regionId, err := ur.GetRegionIdByName(ctx, u.Region)
	if err != nil {
		return fmt.Errorf("GetRegionIdByName: %w", err)
	}
	_, err = ur.pg.Exec(ctx, insertUniversityRequest,
		u.Name, u.Site, u.Prestige, u.Rank,
		u.Quality, u.Scholarship, u.Dormitory,
		u.Labs, u.Sport, regionId)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "", zap.String("Name", u.Name))
		return fmt.Errorf("FillUniversity: %w", err)
	}
	return nil
}

// InsertSpeciality inserts speciality into DB
func (ur *UniversityRepository) InsertSpeciality(ctx context.Context, s *entities.Speciality) error {
	universityId, err := ur.GetUniversityIdByName(ctx, s.UniversityName)
	if err != nil {
		return fmt.Errorf("GetUniversityIdByName: %w", err)
	}
	_, err = ur.pg.Exec(ctx, insertSpecialityRequest,
		universityId, s.Name, s.BudgetPoints, s.ContractPoints, s.Cost)
	if err != nil {
		return fmt.Errorf("FillSpeciality: %w", err)
	}
	return nil
}
