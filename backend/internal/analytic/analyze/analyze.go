package analyze

import (
	"University-Selection-Service/internal/entities"
	"slices"
)

type Analyser struct{}

// GetCriteriaWeights creates a comparison matrix and calculates the weights of the criteria
func (a *Analyser) GetCriteriaWeights(crt *entities.Criteria, cmp *entities.Comparisons) *entities.Criteria {
	var matrix [4][4]float64

	var col0 float64 = 0
	var col1 float64 = 0
	var col2 float64 = 0
	var col3 float64 = 0

	matrix[0][1] = float64(cmp.RatingToPrestige)
	matrix[0][2] = float64(cmp.RatingToEducationQuality)
	matrix[0][3] = float64(cmp.RatingToScholarshipPrograms)
	matrix[1][2] = float64(cmp.PrestigeToEducationQuality)
	matrix[1][3] = float64(cmp.PrestigeToScholarshipPrograms)
	matrix[2][3] = float64(cmp.EducationQualityToScholarshipPrograms)

	// fill matrix with comparisons
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i == j {
				matrix[i][j] = 1
			} else if i > j {
				matrix[i][j] = 1.0 / matrix[j][i]
			}
			switch j {
			case 0:
				col0 += matrix[i][j]
			case 1:
				col1 += matrix[i][j]
			case 2:
				col2 += matrix[i][j]
			case 3:
				col3 += matrix[i][j]
			}
		}
	}

	// normalization of columns
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			switch j {
			case 0:
				matrix[i][j] = matrix[i][j] / col0
			case 1:
				matrix[i][j] = matrix[i][j] / col1
			case 2:
				matrix[i][j] = matrix[i][j] / col2
			case 3:
				matrix[i][j] = matrix[i][j] / col3
			}
		}
	}

	// get average in row
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			switch i {
			case 0:
				crt.LocalUniversityRating += matrix[i][j]
			case 1:
				crt.Prestige += matrix[i][j]
			case 2:
				crt.EducationQuality += matrix[i][j]
			case 3:
				crt.ScholarshipPrograms += matrix[i][j]
			}
		}
	}
	crt.LocalUniversityRating = crt.LocalUniversityRating / 4.0
	crt.Prestige = crt.Prestige / 4.0
	crt.EducationQuality = crt.EducationQuality / 4.0
	crt.ScholarshipPrograms = crt.ScholarshipPrograms / 4.0

	return crt
}

// Analyze analyzes universities for compliance based on paired comparison criteria
func (a *Analyser) Analyze(universities []*entities.University, cmp *entities.Comparisons, rankSum float64, prestigeSum, educationQualitySum, scholarshipProgramsSum int) ([]*entities.University, error) {
	criteria := &entities.Criteria{}
	criteria = a.GetCriteriaWeights(criteria, cmp)

	for _, univ := range universities {
		univ.Relevancy = univ.Relevancy/rankSum*criteria.LocalUniversityRating +
			float64(univ.Prestige)/float64(prestigeSum)*criteria.Prestige +
			float64(univ.Quality)/float64(educationQualitySum)*criteria.EducationQuality +
			float64(univ.Scholarship)/float64(scholarshipProgramsSum)*criteria.ScholarshipPrograms
	}
	slices.SortFunc(universities, func(u1, u2 *entities.University) int {
		return int(u1.Relevancy - u2.Relevancy)
	})
	return universities, nil
}
