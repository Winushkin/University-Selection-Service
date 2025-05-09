package entities

type University struct {
	Id             int
	Prestige       int
	Name           string
	Site           string
	Rank           float64
	Quality        int
	ContractPoints float64
	BudgetPoints   float64
	Cost           int
	Dormitory      bool
	Labs           bool
	Sport          bool
	Scholarship    int
	Region         string
}
