package university

import (
	"University-Selection-Service/internal/entities"
	"fmt"
	"math/rand/v2"
	"strconv"
)

func FillUniversityData(university *entities.University, num int) {
	rank, err := strconv.ParseFloat(fmt.Sprintf("%.1f", rand.Float64()*2+3), 64)
	if err != nil {
		fmt.Println(err)
	}
	university.Id = num
	university.Prestige = num
	university.Rank = rank
	university.Quality = rand.IntN(5) + 5
	university.Cost = rand.IntN(200000) + 100000
	university.Scholarship = rand.IntN(8000) + 2000
}

func CreateSpecialitiesByUniversity(specialitiesSet []string,
	university *entities.University) []*entities.Speciality {
	specialities := make([]*entities.Speciality, 0)
	for i := 0; i < rand.IntN(len(specialitiesSet)); i++ {

		cost := university.Cost - 30000 + rand.IntN(60000)

		points := rand.IntN(20) - 10
		budgetPoints, contractPoints := 0, 0
		if university.BudgetPoints != 0 {
			budgetPoints = points
		}
		if university.BudgetPoints != 0 {
			contractPoints = points
		}

		speciality := &entities.Speciality{
			UniversityName: university.Name,
			Name:           specialitiesSet[rand.IntN(len(specialitiesSet))],
			BudgetPoints:   university.BudgetPoints + budgetPoints,
			ContractPoints: university.ContractPoints + contractPoints,
			Cost:           cost - cost%1000,
		}

		if speciality.BudgetPoints > 100 {
			speciality.BudgetPoints = 100
		}
		if speciality.ContractPoints > 100 {
			speciality.ContractPoints = 100
		}

		specialities = append(specialities, speciality)
	}
	return specialities
}
