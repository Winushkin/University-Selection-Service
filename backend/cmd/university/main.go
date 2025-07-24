package main

import (
	"University-Selection-Service/internal/config"
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/internal/parser"
	"University-Selection-Service/internal/repositories"
	"University-Selection-Service/internal/university"
	"University-Selection-Service/pkg/logger"
	"context"
	"go.uber.org/zap"
)

// main starts university service
func main() {
	ctx := context.Background()
	ctx, _ = logger.NewLogger(ctx)
	log := logger.GetLoggerFromCtx(ctx)

	cfg, err := config.NewUniversityConfig()
	if err != nil {
		log.Error(ctx, "cannot load university config", zap.Error(err))
		return
	}

	universities, specialitiesSet, regions, err := parser.ParseUniversities(ctx, cfg.DatasetPath)
	if err != nil {
		log.Error(ctx, "cannot parse universities", zap.Error(err))
	}

	for i := 0; i < len(universities); i++ {
		university.FillUniversityData(&universities[i], i+1)
	}

	specialities := make([]*entities.Speciality, 0)
	for _, univer := range universities {
		specs := university.CreateSpecialitiesByUniversity(specialitiesSet, &univer)
		specialities = append(specialities, specs...)
	}

	var repository repositories.UniversityRepoInterface
	repository, err = repositories.NewUniversityRepository(ctx, cfg.Postgres)
	if err != nil {
		log.Error(ctx, "failed to connect to university repository", zap.Error(err))
		return
	}

	err = repository.FillRegions(ctx, regions)
	if err != nil {
		log.Error(ctx, "failed to fill regions", zap.Error(err))
	}

	for _, univer := range universities {
		err = repository.InsertUniversity(ctx, &univer)
		if err != nil {
			log.Error(ctx, "failed to fill university", zap.Error(err))
		}
	}

	for _, speciality := range specialities {
		err = repository.InsertSpeciality(ctx, speciality)
		if err != nil {
			log.Error(ctx, "failed to insert speciality", zap.Error(err))
		}
	}
}
