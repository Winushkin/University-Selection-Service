package parser

import (
	"University-Selection-Service/integrations/config"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

const added_URL = "&ptype=0&glist=0&vuz-abiturients-paid-order=ge&vuz-abiturients-paid-val=&price-order=ge&price-val="

var (
	regions = [84]string{"Алтайский край", "Амурская область", "Архангельская область",
		"Астраханская область", "Белгородская область", "Брянская область",
		"Владимирская область", "Волгоградская область", "Вологодская область",
		"Воронежская область", "Еврейская автономная область", "Забайкальский край",
		"Ивановская область", "Иркутская область", "Кабардино-Балкарская Республика",
		"Калининградская область", "Калужская область", "Камчатский край",
		"Карачаево-Черкесская Республика", "Кемеровская область", "Кировская область",
		"Костромская область", "Краснодарский край", "Красноярский край",
		"Курганская область", "Курская область", "Ленинградская область",
		"Липецкая область", "Магаданская область", "Москва и Московская область",
		"Мурманская область", "Ненецкий автономный округ", "Нижегородская область",
		"Новгородская область", "Новосибирская область", "Омская область",
		"Оренбургская область", "Орловская область", "Пензенская область",
		"Пермский край", "Приморский край", "Псковская область",
		"Республика Адыгея", "Республика Алтай", "Республика Башкортостан",
		"Республика Бурятия", "Республика Дагестан", "Республика Ингушетия",
		"Республика Калмыкия", "Республика Карелия", "Республика Коми",
		"Республика Крым", "Республика Марий Эл", "Республика Мордовия",
		"Республика Саха (Якутия)", "Республика Северная Осетия — Алания", "Республика Татарстан",
		"Республика Тыва", "Республика Хакасия", "Ростовская область",
		"Рязанская область", "Самарская область", "Санкт-Петербург",
		"Саратовская область", "Сахалинская область", "Свердловская область",
		"Севастополь", "Смоленская область", "Ставропольский край",
		"Тамбовская область", "Тверская область", "Томская область",
		"Тульская область", "Тюменская область", "Удмуртская Республика",
		"Ульяновская область", "Хабаровский край", "Ханты-Мансийский автономный округ — Югра",
		"Челябинская область", "Чеченская Республика", "Чувашская Республика",
		"Чукотский автономный округ", "Ямало-Ненецкий автономный округ", "Ярославская область"}
)

func Parse(url string) (*[][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != 200 {
		fmt.Printf("Ошибка: статус код %d\n", resp.StatusCode)
		return nil, err
	}
	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка при парсинге HTML: %v\n", err)
		return nil, err
	}

	// Find the table and iterate over rows
	var result [][]string
	doc.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		// Extract data from each column
		broadenedGroup := strings.TrimSpace(row.Find("td").Eq(0).Text())
		university := strings.TrimSpace(row.Find("td").Eq(1).Text())
		avgEGEScore := strings.TrimSpace(row.Find("td").Eq(2).Text())

		Group := []string{broadenedGroup, university, avgEGEScore}
		result = append(result, Group)

	})
	return &result, nil
}

func ParseAllData(cfg config.IntegrationConfig) (*[][]string, *[][]string, error) {
	var url string
	var budget [][]string
	var contract [][]string

	for _, region := range regions {
		region = strings.ReplaceAll(region, " ", "+")

		url = fmt.Sprintf("%s%s%s", cfg.BudgetURL, region, added_URL)
		subBudget, err := Parse(url)
		if err != nil {
			return nil, nil, err
		}

		for _, sub := range *subBudget {
			budget = append(budget, sub)
		}

		url = fmt.Sprintf("%s%s%s", cfg.ContractURL, region, added_URL)
		subContract, err := Parse(url)
		if err != nil {
			return nil, nil, err
		}

		for _, sub := range *subContract {
			contract = append(contract, sub)
		}
	}
	return &budget, &contract, nil

}
