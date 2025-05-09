package parser

import (
	"University-Selection-Service/internal/entities"
	"context"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

func Set(collection []string) ([]string, error) {
	CollectionMap := make(map[string]bool)
	CollectionSet := make([]string, 0)

	for _, item := range collection {
		CollectionMap[item] = true
	}

	for item, _ := range CollectionMap {
		CollectionSet = append(CollectionSet, item)
	}
	return CollectionSet, nil
}

func ParseUniversities(ctx context.Context, path string) ([]entities.University, []string, []string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)

		}
	}()

	cols, err := f.GetCols("Universities")
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}

	regions, err := Set(cols[1][1:])
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}

	specialities, err := Set(cols[8][1:])
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}

	rows, err := f.GetRows("Universities")
	rows = rows[1:]
	universities := make([]entities.University, 0)
	var university entities.University

	for _, row := range rows {
		budgetPoints, _ := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
		contractPoints, _ := strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
		if budgetPoints == 0 && contractPoints == 0 {
			continue
		}

		university = entities.University{
			Name:           strings.TrimSpace(row[0]),
			Site:           strings.TrimSpace(row[2]),
			ContractPoints: contractPoints,
			BudgetPoints:   budgetPoints,
			Dormitory:      strings.TrimSpace(row[5]) != "0",
			Labs:           strings.TrimSpace(row[7]) != "0",
			Sport:          strings.TrimSpace(row[6]) != "0",
			Region:         strings.TrimSpace(row[1]),
		}
		universities = append(universities, university)

	}

	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}

	return universities, specialities, regions, nil

}
