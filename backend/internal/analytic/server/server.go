package analytic

import (
	"University-Selection-Service/internal/analytic/analyze"
	"University-Selection-Service/internal/analytic/repository"
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/pkg/api"
	"context"
	"fmt"
)

type Server struct {
	api.AnalyticServer
	userCli      api.UserServiceClient
	RepInterface repository.AnalyticRepositoryInterface
}

// New creates new Analytic Server
func New(client api.UserServiceClient, r interface {
	GetUniversitiesBySpeciality(ctx context.Context, specialityName string) ([]*entities.University, error)
}) (*Server, error) {
	return &Server{
		userCli:      client,
		RepInterface: r,
	}, nil
}

// Analyze handles /api/analytic/analyze request
func (s *Server) Analyze(ctx context.Context, request *api.AnalyzeRequest) (*api.AnalyzeResponse, error) {
	id, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, fmt.Errorf("id is not int")
	}
	user, err := s.userCli.ProfileDataForAnalytic(ctx, &api.ProfileDataForAnalyticRequest{
		Id: int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("ProfileDataForAnalytic: %w", err)
	}

	univResult := &api.AnalyzeResponse{
		Speciality: user.Speciality,
	}

	universities, rankSum, prestigeSum, educationQualitySum, scholarshipProgramsSum, err := s.FilterUniversities(ctx, user, int(request.EducationCost), request.Dormitory,
		request.SportsInfrastructure, request.ScientificLabs)
	if err != nil {
		return nil, fmt.Errorf("FilterUniversities: %w", err)
	}

	a := &analyze.Analyser{}

	universities, err = a.Analyze(universities, &entities.Comparisons{
		RatingToPrestige:                      int(request.RatingToPrestige),
		RatingToEducationQuality:              int(request.RatingToEducationQuality),
		RatingToScholarshipPrograms:           int(request.RatingToScholarshipPrograms),
		PrestigeToEducationQuality:            int(request.PrestigeToEducationQuality),
		PrestigeToScholarshipPrograms:         int(request.PrestigeToScholarshipPrograms),
		EducationQualityToScholarshipPrograms: int(request.EducationQualityToScholarshipPrograms),
	}, rankSum, prestigeSum, educationQualitySum, scholarshipProgramsSum)

	for _, univ := range universities {
		u := &api.University{
			Name:           univ.Name,
			Region:         univ.Region,
			BudgetPoints:   int32(univ.BudgetPoints),
			ContractPoints: int32(univ.ContractPoints),
			Cost:           int64(univ.Cost),
			Prestige:       int32(univ.Prestige),
			Rank:           float32(univ.Rank),
			Quality:        int32(univ.Quality),
			Dormitory:      univ.Dormitory,
			Labs:           univ.Labs,
			Sport:          univ.Sport,
			Scholarship:    int32(univ.Scholarship),
			Relevancy:      univ.Relevancy,
			Site:           univ.Site,
		}
		univResult.Universities = append(univResult.Universities, u)
	}

	return univResult, nil
}

// FilterUniversities filters universities with params
func (s *Server) FilterUniversities(ctx context.Context, user *api.ProfileDataForAnalyticResponse, cost int, dormitory bool, sport bool, labs bool) ([]*entities.University, float64, int, int, int, error) {
	universities, err := s.RepInterface.GetUniversitiesBySpeciality(ctx, user.Speciality)
	rankSum, prestigeSum, educationQualitySum, scholarshipProgramsSum := 0.0, 0, 0, 0
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("GetUniversitiesBySpeciality: %w", err)
	}

	filteredUniversities := make([]*entities.University, 0)
	for _, univ := range universities {
		if univ.Region != user.Town {
			continue
		}

		if user.Financing == "Бюджет" {
			if int(user.Ege) < univ.BudgetPoints {
				continue
			}
		} else {
			if int(user.Ege) < univ.ContractPoints {
				continue
			}
			if univ.Cost > cost {
				continue
			}
		}

		if dormitory == true && univ.Dormitory == false {
			continue
		}
		if sport == true && univ.Sport == false {
			continue
		}
		if labs == true && univ.Labs == false {
			continue
		}

		filteredUniversities = append(filteredUniversities, univ)
		rankSum += univ.Rank
		prestigeSum += univ.Prestige
		educationQualitySum += univ.Quality
		scholarshipProgramsSum += univ.Scholarship
	}
	return filteredUniversities, rankSum, prestigeSum, educationQualitySum, scholarshipProgramsSum, nil
}
