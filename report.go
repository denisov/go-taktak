package main

import (
	"fmt"
	"github.com/goodsign/monday"
	"github.com/tealeg/xlsx"
	"gopkg.in/gomail.v2"
	"log"
	"sort"
	"strconv"
	"time"
)

// writeReportFile создаёт файл с данными отчёта
func writeReportFile(problems []problem, fileName string) {
	log.Println("Создаём отчёт")

	sort.Slice(problems, func(i, j int) bool {
		return problems[i].solutionDate.Before(problems[j].solutionDate)
	})

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		log.Fatalln(err)
	}
	addHeader(sheet)

	idx := 1
	for _, problemItem := range problems {
		row := sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(idx)
		row.AddCell().Value = problemItem.solutionDate.Format("02.01.2006")
		row.AddCell().Value = problemItem.title
		row.AddCell().Value = problemItem.url

		idx++
	}

	sheet.Col(0).Width = 3
	sheet.Col(1).Width = 13
	sheet.Col(2).Width = 85
	sheet.Col(3).Width = 22

	err = file.Save(fileName)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

// getReportName возвращает имя файла на основе id пользователя и даты отчёта
func getReportName(userId, month, year int) string {
	name := strconv.Itoa(userId)
	if surname, ok := cfg.Users[userId]; ok {
		name = surname
	}
	return fmt.Sprintf(
		"Отчёт_taktaktak_%s_%s.xlsx",
		monday.Format(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc), "January_2006", monday.LocaleRuRU),
		name,
	)
}

// addHeader добавляет строку - заголовок таблицы отчёта
func addHeader(sheet *xlsx.Sheet) {
	row := sheet.AddRow()
	titles := []string{"№", "Даты консультации", "Вопрос", "Краткий комментарий (ответ)"}
	for _, titleItem := range titles {
		row.AddCell().Value = titleItem
	}
}

// sendReport отправляет файл отчёта на email
func sendReport(fileName, subject string) {
	log.Println("Отправляем на email")

	message := gomail.NewMessage()
	message.SetHeader("From", cfg.Email.From)
	message.SetHeader("To", cfg.Email.To)
	message.SetHeader("Subject", subject)

	message.Attach(fileName)

	port, err := strconv.Atoi(cfg.Email.Port)
	if err != nil {
		log.Fatalln(err)
	}

	dialer := gomail.NewDialer(
		cfg.Email.Host,
		port,
		cfg.Email.From,
		cfg.Email.Password,
	)

	if err := dialer.DialAndSend(message); err != nil {
		log.Fatalln(err)
	}
}
