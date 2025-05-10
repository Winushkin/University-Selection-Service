package entities

type University struct {
	Id             int
	Prestige       int
	Name           string
	Site           string
	Rank           float64
	Quality        int
	ContractPoints int
	BudgetPoints   int
	Cost           int
	Dormitory      bool
	Labs           bool
	Sport          bool
	Scholarship    int
	Region         string
	Relevancy      float64
}
