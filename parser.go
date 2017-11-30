package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"crypto/tls"
	"github.com/PuerkitoBio/goquery"
	"github.com/goodsign/monday"
	"net/http"
)

// возвращает статистику проблем, решённых пользователем
func getUserStat(userId, month, year, monthDeep int) (problems []problem) {
	maxDate := time.Date(year, time.Month(month-monthDeep+1), 1, 0, 0, 0, 0, loc)

	log.Printf(
		"Ищем решения проблем для пользователя %d за %s. До %s\n\n",
		userId,
		monday.Format(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc), "January 2006", monday.LocaleRuRU),
		monday.Format(maxDate, monday.DefaultFormatRuRULong, monday.LocaleRuRU),
	)

	page := 1
	for {
		problemsSlice := getProblems(userId, page)

		if len(problemsSlice) == 0 {
			log.Print("Проблемы закончились")
			return
		}

		for _, problemItem := range problemsSlice {

			solutionDate := getSolutionDate(problemItem.url, userId)

			if solutionDate.IsZero() {
				log.Printf(
					"Проблема: %-55s (%s)   нулевая дата, пропускаем. Коммент удалён?",
					truncate(problemItem.title, 50, ".."),
					monday.Format(problemItem.problemDate, monday.DefaultFormatRuRUMedium, monday.LocaleRuRU),
				)
				continue
			}
			if solutionDate.Before(maxDate) {
				log.Printf("Время решения превысило максимальную дату. %s => %v", problemItem.title, solutionDate)
				return
			}

			var solutionString string
			if solutionDate.Year() == year && int(solutionDate.Month()) == month {
				problemItem.solutionDate = solutionDate
				solutionString = "Решение: " + monday.Format(solutionDate, monday.DefaultFormatRuRUMedium, monday.LocaleRuRU)
				problems = append(problems, problemItem)
			}

			log.Printf(
				"Проблема: %-55s (%s)   %s",
				truncate(problemItem.title, 50, ".."),
				monday.Format(problemItem.problemDate, monday.DefaultFormatRuRUMedium, monday.LocaleRuRU),
				solutionString,
			)

		}
		page++
	}
}

// getProblems возвращает проблемы, решённые прользователем
func getProblems(userId, page int) []problem {
	url := fmt.Sprintf("https://taktaktak.ru/person/%d/answers?page=%d&ajax=2", userId, page)
	resp := getResponse(url)
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatalf("Не могу создать документ для парсинга URL: %s. %s", url, err)
	}

	var problems []problem
	doc.Find(".item").Each(func(i int, item *goquery.Selection) {
		title := item.Find("h3 a").First()
		problemPageUrl, exists := title.Attr("href")
		if !exists {
			html, _ := title.Html()
			log.Printf("Не могу найти аттрибут 'href' в '%s'", html)
			return
		}

		dateText := item.Find("span.date span").Text()
		problemDate := time.Time{}
		if dateText == "" {
			log.Print("Не могу найти дату")
		} else {
			problemDate, err = parseDate(dateText)
			if err != nil {
				log.Println(err)
			}
		}

		problemItem := &problem{
			title:       strings.Trim(title.Text(), "\n "),
			url:         fmt.Sprintf("https://taktaktak.ru%s", problemPageUrl),
			problemDate: problemDate,
		}
		problems = append(problems, *problemItem)
	})

	return problems
}

// getSolutionDate возвращает дату решения проблемы. Дата долждна быть всегда, если нет - вернёт нулевое время
func getSolutionDate(url string, userId int) time.Time {
	resp := getResponse(url)
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatalf("Не могу создать документ для парсинга URL: %s. %s", url, err)
	}

	time.Sleep(2 * time.Second)
	if err != nil {
		log.Fatalf("Не могу открыть URL: %s. %s", url, err)
	}

	var solutionDate time.Time
	selector := "div.answer a[href=\"/person/" + strconv.Itoa(userId) + "\"]"
	doc.Find(selector).EachWithBreak(func(i int, selection *goquery.Selection) bool {
		solutionDateText := selection.Parent().Parent().Find(".date span").First().Text()
		// у удалённых комментов нет даты, продолжаем искать дальше
		if solutionDateText == "" {
			return true
		}

		solutionDate, err = parseDate(solutionDateText)
		if err != nil {
			log.Fatal(err)
		}
		return false // остановить перебор
	})

	return solutionDate
}

func getResponse(urlToParse string) *http.Response {

	// FIXME не надо каждый раз создавать
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(urlToParse)
	if err != nil {
		log.Fatalf("Не могу открыть URL: %s. %s", urlToParse, err)
	}

	return resp
}
