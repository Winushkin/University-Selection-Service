package parser

import (
	"University-Selection-Service/internal/config"
	"University-Selection-Service/internal/entities"
	"University-Selection-Service/pkg/logger"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const addedUrl = "&ptype=0&glist=0&vuz-abiturients-paid-order=ge&vuz-abiturients-paid-val=&price-order=ge&price-val="

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

func parse(ctx context.Context, url string) (*[]*entities.University, error) {
	log := logger.GetLoggerFromCtx(ctx)
	resp, err := http.Get(url)
	if err != nil {
		log.Error(ctx, "failed to fetch url", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != 200 {
		log.Error(ctx, "ERROR, status code: ", zap.Int("status", resp.StatusCode))
		return nil, err
	}
	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error(ctx, "failed to parse url", zap.Error(err))
		return nil, err
	}

	// Find the table and iterate over rows
	var result []*entities.University
	var univer *entities.University
	doc.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		// Extract data from each column
		speciality := strings.TrimSpace(row.Find("td").Eq(0).Text())
		university := strings.TrimSpace(row.Find("td").Eq(1).Text())
		avgEGEScore, _ := strconv.ParseFloat(strings.TrimSpace(row.Find("td").Eq(2).Text()), 64)

		univer = &entities.University{
			Name:        university,
			Rank:        rand.Float64()*2 + 3,
			Speciality:  speciality,
			Quality:     rand.IntN(5) + 5,
			Points:      avgEGEScore,
			Cost:        rand.IntN(200000) + 100000,
			Dormitory:   rand.IntN(5) != 0,
			Labs:        rand.IntN(3) != 0,
			Sport:       rand.IntN(5) != 0,
			Scholarship: rand.IntN(8000) + 2000,
		}
		result = append(result, univer)
	})

	result = result[6:]
	return &result, nil
}

func ParseAllData(ctx context.Context, cfg *config.UniversityConfig) (*[]*entities.University,
	*[]*entities.University, error) {

	var url string
	var budget []*entities.University
	var contract []*entities.University

	log := logger.GetLoggerFromCtx(ctx)

	wg := &sync.WaitGroup{}
	var m sync.Mutex
	var formatedRegion string
	for _, region := range regions {
		wg.Add(1)
		formatedRegion = strings.ReplaceAll(region, " ", "+")

		go func(m *sync.Mutex, formatedRegion string, cfg *config.UniversityConfig) {
			defer wg.Done()
			url = fmt.Sprintf("%s%s%s", cfg.BudgetURL, formatedRegion, addedUrl)
			subBudget, err := parse(ctx, url)
			if err != nil || subBudget == nil {
				log.Error(ctx, "failed to parse budget", zap.Error(err))
				return
			}
			m.Lock()
			for _, univer := range *subBudget {
				univer.Region = region
				budget = append(budget, univer)
			}
			m.Unlock()

			url = fmt.Sprintf("%s%s%s", cfg.ContractURL, formatedRegion, addedUrl)
			subContract, err := parse(ctx, url)
			if err != nil || subContract == nil {
				log.Error(ctx, "failed to parse contract", zap.Error(err))
				return
			}

			m.Lock()
			for _, univer := range *subBudget {
				univer.Region = region
				contract = append(contract, univer)
			}
			m.Unlock()

		}(&m, formatedRegion, cfg)
	}
	wg.Wait()
	return &budget, &contract, nil

}
