package main

import (
	"fmt"
	"github.com/goodsign/monday"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var loc *time.Location

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		log.Fatal("Не могу загрузить location ", err)
	}
}

func parseDate(taktakDate string) (time.Time, error) {

	taktakDate = strings.Trim(taktakDate, " ")
	// дата без года - добавляем год
	re := regexp.MustCompile(`(\d+ [а-я]+), (\d+)`)
	taktakDate = re.ReplaceAllString(taktakDate, fmt.Sprintf("$1 %d, $2", time.Now().Year()))

	// добавляем текущую дату
	re = regexp.MustCompile(`(сегодня), (\d+:\d+)`)
	taktakDate = re.ReplaceAllString(taktakDate, time.Now().Format("2 January 2006")+", $2")

	// добавляем вчерашнюю дату
	re = regexp.MustCompile(`(вчера), (\d+:\d+)`)
	yesterday := time.Now().AddDate(0, 0, -1)
	taktakDate = re.ReplaceAllString(taktakDate, yesterday.Format("2 January 2006")+", $2")

	if taktakDate == "1 час назад" {
		return time.Now().Add(-1 * time.Hour), nil
	}

	if taktakDate == "2 часа назад" {
		return time.Now().Add(-1 * time.Hour), nil
	}

	parsedDate, err := monday.ParseInLocation(
		"2 January 2006, 15:04",
		taktakDate,
		loc,
		monday.LocaleRuRU,
	)
	if err != nil {
		return time.Time{}, fmt.Errorf("Не могу распарсить дату %s. %s", taktakDate, err)
	}

	return parsedDate, nil
}

// truncate обрезает строку до заданной длины и дополняет её многоточсием
func truncate(str string, maxLength int, ending string) string {
	strRunes := []rune(str)
	if len(strRunes) <= maxLength {
		return str
	}

	return string(strRunes[0:maxLength-len(ending)]) + ending
}

// ensureDir проверяет если ли заданная папка, создаёт если нет
func ensureDir(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0775)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	if !stat.IsDir() {
		return fmt.Errorf("%s существует, но не явлется директорией", dir)
	}

	return nil
}
