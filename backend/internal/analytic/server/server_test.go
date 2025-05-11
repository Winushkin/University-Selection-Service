package analytic

import (
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/pkg/api"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockAnalyticRepository implement AnalyticRepositoryInterface
type MockAnalyticRepository struct {
	mock.Mock
}

func (m *MockAnalyticRepository) GetUniversitiesBySpeciality(ctx context.Context, specialityName string) ([]*entities.University, error) {
	args := m.Called(ctx, specialityName)
	return args.Get(0).([]*entities.University), args.Error(1)
}

func TestFilterUniversities_BudgetFilter(t *testing.T) {
	repo := new(MockAnalyticRepository)
	s := &Server{RepInterface: repo}

	user := &api.ProfileDataForAnalyticResponse{
		Speciality: "CS",
		Financing:  "Бюджет",
		Ege:        250,
	}

	testUniv := []*entities.University{
		{BudgetPoints: 200, ContractPoints: 150},
		{BudgetPoints: 300, ContractPoints: 250},
	}

	repo.On("GetUniversitiesBySpeciality", mock.Anything, "CS").Return(testUniv, nil)

	filtered, _, _, _, _, err := s.FilterUniversities(
		context.Background(),
		user,
		0,
		false,
		false,
		false,
	)

	assert.NoError(t, err)
	assert.Len(t, filtered, 1)
	assert.Equal(t, 200, filtered[0].BudgetPoints)
}

func TestFilterUniversities_ContractFilter(t *testing.T) {
	repo := new(MockAnalyticRepository)
	s := &Server{RepInterface: repo}

	user := &api.ProfileDataForAnalyticResponse{
		Speciality: "CS",
		Financing:  "Контракт",
		Ege:        200,
	}

	testUniv := []*entities.University{
		{ContractPoints: 150, Cost: 100000},
		{ContractPoints: 250, Cost: 50000},
		{ContractPoints: 150, Cost: 200000},
	}

	repo.On("GetUniversitiesBySpeciality", mock.Anything, "CS").Return(testUniv, nil)

	filtered, _, _, _, _, err := s.FilterUniversities(
		context.Background(),
		user,
		150000,
		false,
		false,
		false,
	)

	assert.NoError(t, err)
	assert.Len(t, filtered, 1)
	assert.Equal(t, 100000, filtered[0].Cost)
}

func TestFilterUniversities_FacilityFilters(t *testing.T) {
	repo := new(MockAnalyticRepository)
	s := &Server{RepInterface: repo}

	user := &api.ProfileDataForAnalyticResponse{
		Speciality: "CS",
		Financing:  "Бюджет",
		Ege:        300,
	}

	testUniv := []*entities.University{
		{Dormitory: true, Sport: true, Labs: true},
		{Dormitory: false, Sport: true, Labs: true},
	}

	repo.On("GetUniversitiesBySpeciality", mock.Anything, "CS").Return(testUniv, nil)

	filtered, _, _, _, _, err := s.FilterUniversities(
		context.Background(),
		user,
		0,
		true,
		true,
		true,
	)

	assert.NoError(t, err)
	assert.Len(t, filtered, 1)
	assert.True(t, filtered[0].Dormitory)
}
