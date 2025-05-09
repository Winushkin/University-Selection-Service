package repositories

import (
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/pkg/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	getUniversitiesBySpecialityRequest = "SELECT u.*, r.name FROM universities.universities u JOIN universities.regions r ON u.region_id = r.id JOIN universities.specialities s ON u.id = s.university_id WHERE s.name = $1"
)

type AnalyticRepositoryInterface interface {
	GetUniversitiesBySpeciality(ctx context.Context, specialityName string)
}

type AnalyticRepository struct {
	pg *pgxpool.Pool
}

func NewAnalyticRepository(ctx context.Context, cfg postgres.Config) (*AnalyticRepository, error) {
	pool, err := postgres.New(ctx, cfg, "universities")
	if err != nil {
		return nil, fmt.Errorf("NewUniversityRepository: failed to connect to users postgres: %w", err)
	}
	return &AnalyticRepository{pg: pool}, nil
}

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
		err = queryRows.Scan(&u.Id, &u.Name, &u.Prestige, &u.Rank, &u.Quality, &u.Scholarship,
			&u.Dormitory, &u.Labs, &u.Sport, &region, &u.Region)
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
