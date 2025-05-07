package main

import (
	"University-Selection-Service/internal/config"
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/internal/parser"
	"University-Selection-Service/internal/repositories"
	"University-Selection-Service/pkg/logger"
	"University-Selection-Service/pkg/postgres"
	"context"
	"go.uber.org/zap"
	"sync"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger.NewLogger(ctx)
	log := logger.GetLoggerFromCtx(ctx)

	cfg, err := config.NewUniversityConfig()
	if err != nil {
		log.Error(ctx, "cannot load university config", zap.Error(err))
		return
	}

	budget, contract, err := parser.ParseAllData(ctx, cfg)

	if err != nil || budget == nil || contract == nil {
		log.Error(ctx, "failed to parse all data", zap.Error(err))
		return
	}

	db, err := postgres.New(ctx, cfg.Postgres, "universities")
	if err != nil {
		log.Error(ctx, "failed to connect to university postgres", zap.Error(err))
		return
	}
	defer db.Close()

	repository, err := repositories.NewUniversityRepository(ctx, cfg.Postgres)
	if err != nil {
		log.Error(ctx, "failed to connect to university repository", zap.Error(err))
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func(ctx context.Context, ctrc *[]*entities.University, rep *repositories.UniversityRepository) {
		defer wg.Done()
		prestige := 0
		for _, university := range *ctrc {
			prestige++
			university.Prestige = prestige
			err = rep.FillContract(ctx, university)
			if err != nil {
				log.Error(ctx, "failed to fill contract", zap.Error(err))
				return
			}
		}
	}(ctx, contract, repository)

	go func(ctx context.Context, bdg *[]*entities.University, rep *repositories.UniversityRepository) {
		defer wg.Done()
		prestige := 0
		for _, university := range *bdg {
			prestige++
			university.Prestige = prestige
			err = rep.FillBudget(ctx, university)
			if err != nil {
				log.Error(ctx, "failed to fill budget", zap.Error(err))
				return
			}
		}
	}(ctx, budget, repository)

	wg.Wait()

}
