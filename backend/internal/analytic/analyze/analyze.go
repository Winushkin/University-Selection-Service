package analyze

import "University-Selection-Service/internal/entities"

type Analyser struct{}

func (a *Analyser) GetCriteriaWeights(crt *entities.Criteria, cmp *entities.Comparisons) (*entities.Criteria, error) {
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

	return crt, nil
}
