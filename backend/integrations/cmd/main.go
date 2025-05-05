package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func main() {
	// URL of the HSE rating page
	url := "https://ege.hse.ru/rating/2024/98633979/all/?rlist=Алтайский+край&ptype=0&glist=0&vuz-abiturients-paid" +
		"-order" +
		"=ge&vuz" +
		"-abiturients-paid-val=&price-order=ge&price-val="

	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != 200 {
		fmt.Printf("Ошибка: статус код %d\n", resp.StatusCode)
		return
	}

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка при парсинге HTML: %v\n", err)
		return
	}

	// Find the table and iterate over rows
	doc.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		// Extract data from each column
		broadenedGroup := strings.TrimSpace(row.Find("td").Eq(0).Text())
		university := strings.TrimSpace(row.Find("td").Eq(1).Text())
		avgEGEScore := strings.TrimSpace(row.Find("td").Eq(2).Text())

		// Print or store the extracted data
		fmt.Printf("Укрупнённая группа: %s\n", broadenedGroup)
		fmt.Printf("Вуз: %s\n", university)
		fmt.Printf("Средний балл ЕГЭ: %s\n", avgEGEScore)
		fmt.Println("---")
	})
}
