package analyze

import (
	"University-Selection-Service/internal/entities"
	"math"
	"testing"
)

func TestGetCriteriaWeights(t *testing.T) {
	tests := []struct {
		name        string
		comparisons entities.Comparisons
		expected    entities.Criteria
	}{
		{
			name: "All comparisons equal to 1",
			comparisons: entities.Comparisons{
				RatingToPrestige:                      1,
				RatingToEducationQuality:              1,
				RatingToScholarshipPrograms:           1,
				PrestigeToEducationQuality:            1,
				PrestigeToScholarshipPrograms:         1,
				EducationQualityToScholarshipPrograms: 1,
			},
			expected: entities.Criteria{
				LocalUniversityRating: 0.25,
				Prestige:              0.25,
				EducationQuality:      0.25,
				ScholarshipPrograms:   0.25,
			},
		},
		{
			name: "Extreme values",
			comparisons: entities.Comparisons{
				RatingToPrestige:                      9,
				RatingToEducationQuality:              1,
				RatingToScholarshipPrograms:           1,
				PrestigeToEducationQuality:            1,
				PrestigeToScholarshipPrograms:         1,
				EducationQualityToScholarshipPrograms: 1,
			},
			expected: entities.Criteria{
				LocalUniversityRating: calculateAverage([]float64{1 / sumCol([]float64{1, 1.0 / 9, 1, 1}), 9 / sumCol([]float64{9, 1, 1, 1}), 1 / sumCol([]float64{1, 1, 1, 1}), 1 / sumCol([]float64{1, 1, 1, 1})}),
				Prestige:              calculateAverage([]float64{(1.0 / 9) / sumCol([]float64{1, 1.0 / 9, 1, 1}), 1 / sumCol([]float64{9, 1, 1, 1}), 1 / sumCol([]float64{1, 1, 1, 1}), 1 / sumCol([]float64{1, 1, 1, 1})}),
				EducationQuality:      calculateAverage([]float64{1 / sumCol([]float64{1, 1.0 / 9, 1, 1}), 1 / sumCol([]float64{9, 1, 1, 1}), 1 / sumCol([]float64{1, 1, 1, 1}), 1 / sumCol([]float64{1, 1, 1, 1})}),
				ScholarshipPrograms:   calculateAverage([]float64{1 / sumCol([]float64{1, 1.0 / 9, 1, 1}), 1 / sumCol([]float64{9, 1, 1, 1}), 1 / sumCol([]float64{1, 1, 1, 1}), 1 / sumCol([]float64{1, 1, 1, 1})}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyser := &Analyser{}
			criteria := &entities.Criteria{}
			result := analyser.GetCriteriaWeights(criteria, &tt.comparisons)
			if !compareCriteria(result, &tt.expected) {
				t.Errorf("Test %s failed: Expected %v, got %v", tt.name, tt.expected, *result)
			}
		})
	}
}

// Helper function to calculate the average of a slice of float64 values
func calculateAverage(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// Helper function to calculate the sum of a slice of float64 values
func sumCol(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum
}

// Helper function to compare two Criteria structs with a tolerance for floating-point errors
func compareCriteria(c1, c2 *entities.Criteria) bool {
	tolerance := 0.0001
	return math.Abs(c1.LocalUniversityRating-c2.LocalUniversityRating) < tolerance &&
		math.Abs(c1.Prestige-c2.Prestige) < tolerance &&
		math.Abs(c1.EducationQuality-c2.EducationQuality) < tolerance &&
		math.Abs(c1.ScholarshipPrograms-c2.ScholarshipPrograms) < tolerance
}
