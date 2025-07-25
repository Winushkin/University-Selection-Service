package repository

import (
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/pkg/postgres"
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed sql/get_universities_by_speciality.sql
	getUniversitiesBySpecialityRequest string

	////go:embed sql/get_regions.sql
	//getRegions string
	//
	////go:embed sql/get_specialities.sql
	//getSpecialities string
)

type AnalyticRepositoryInterface interface {
	GetUniversitiesBySpeciality(ctx context.Context, specialityName string) ([]*entities.University, error)
	//GetRegions(ctx context.Context) ([]string, error)
	//GetSpecialities(ctx context.Context) ([]string, error)
}

type AnalyticRepository struct {
	pg *pgxpool.Pool
}

// NewAnalyticRepository returns new analytic repository with connection to DB
func NewAnalyticRepository(ctx context.Context, cfg postgres.Config) (*AnalyticRepository, error) {
	pool, err := postgres.New(ctx, cfg, "universities")
	if err != nil {
		return nil, fmt.Errorf("NewUniversityRepository: failed to connect to users postgres: %w", err)
	}
	return &AnalyticRepository{pg: pool}, nil
}

// GetUniversitiesBySpeciality returns slice of universities with specific speciality
func (ur *AnalyticRepository) GetUniversitiesBySpeciality(ctx context.Context,
	specialityName string) ([]*entities.University, error) {
	queryRows, err := ur.pg.Query(ctx, getUniversitiesBySpecialityRequest, specialityName)
	if err != nil {
		return nil, fmt.Errorf("GetUniversityBySpeciality: %w", err)
	}
	var result []*entities.University
	for queryRows.Next() {
		var u entities.University
		region := 0
		err = queryRows.Scan(&u.Id, &u.Name, &u.Site, &u.Prestige, &u.Rank, &u.Quality, &u.Scholarship,
			&u.Dormitory, &u.Labs, &u.Sport, &region, &u.Region, &u.Cost, &u.BudgetPoints, &u.ContractPoints)
		if err != nil {
			return nil, fmt.Errorf("GetUniversityBySpeciality Scanning: %w", err)
		}
		result = append(result, &u)
	}

	if err = queryRows.Err(); err != nil {
		return nil, fmt.Errorf("GetUniversityBySpeciality QueryRows: %w", err)
	}
	return result, nil
}

//func (ur *AnalyticRepository) GetRegions(ctx context.Context) ([]string, error) {
//	rows, err := ur.pg.Query(ctx, getRegions)
//	if err != nil {
//		return nil, fmt.Errorf("GetRegions: %w", err)
//	}
//
//	var regions []string
//	for rows.Next() {
//		var region string
//		err = rows.Scan(&region)
//		regions = append(regions, region)
//	}
//
//	return regions, nil
//}
//
//func (ur *AnalyticRepository) GetSpecialities(ctx context.Context) ([]string, error) {
//	rows, err := ur.pg.Query(ctx, getSpecialities)
//	if err != nil {
//		return nil, fmt.Errorf("GetSpecialities: %w", err)
//	}
//	var specialities []string
//	for rows.Next() {
//		var spec string
//		err = rows.Scan(&spec)
//		specialities = append(specialities, spec)
//	}
//	return specialities, nil
//}
