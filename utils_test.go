package main

import (
	"log"
	"testing"
	"time"
)

type ParseDateCase struct {
	taktakDate string
	expected   time.Time
}

func TestParseDate(t *testing.T) {

	var err error
	loc, err = time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		log.Fatal("Не могу загрузить location ", err)
	}

	// простые случаи
	tests := []ParseDateCase{
		{"23 декабря 2016, 10:35", time.Date(2016, 12, 23, 10, 35, 0, 0, loc)},
		{"23 декабря 2016, 7:35", time.Date(2016, 12, 23, 7, 35, 0, 0, loc)},
		{"3 декабря 2016, 7:35", time.Date(2016, 12, 3, 7, 35, 0, 0, loc)},
		{"  3 декабря 2016, 7:35  ", time.Date(2016, 12, 3, 7, 35, 0, 0, loc)},
	}

	tests = append(tests, ParseDateCase{
		"17 апреля, 10:35",
		time.Date(time.Now().Year(), 4, 17, 10, 35, 0, 0, loc),
	})

	today := time.Now()
	tests = append(tests, ParseDateCase{
		"сегодня, 10:35",
		time.Date(today.Year(), today.Month(), today.Day(), 10, 35, 0, 0, loc),
	})

	yesterday := time.Now().AddDate(0, 0, -1)
	tests = append(tests, ParseDateCase{
		"вчера, 10:35",
		time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 10, 35, 0, 0, loc),
	})

	for _, testCase := range tests {
		actual, err := parseDate(testCase.taktakDate)
		if err != nil {
			t.Error(err)
		}
		if actual != testCase.expected {
			t.Errorf(
				"\ntaktakDate: %s\nExpected:   %s\nActual:     %s",
				testCase.taktakDate,
				testCase.expected,
				actual,
			)
		}
	}

	// особо тяжёлые случаи
	actual, err := parseDate("1 час назад")
	if err != nil {
		t.Error(err)
	}
	hourAgo := time.Now().Add(-1 * time.Hour)
	if diff := hourAgo.Sub(actual).Minutes(); diff > 1 {
		t.Error("Не могу разобрать время 1 час назад", diff)
	}

	actual, err = parseDate("2 часа назад")
	if err != nil {
		t.Error(err)
	}
	twoHoursAgo := time.Now().Add(-2 * time.Hour)
	if diff := twoHoursAgo.Sub(actual).Minutes(); diff > 1 {
		t.Error("Не могу разобрать время 2 часа назад", diff)
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		str       string
		maxLength int
		ending    string
		expected  string
	}{
		{"abc", 10, "..", "abc"},
		{"abcde", 3, "..", "a.."},
		{"привет", 5, "..", "при.."},
		{"привет", 500, "..", "привет"},
	}

	for _, testCase := range tests {
		actual := truncate(testCase.str, testCase.maxLength, testCase.ending)

		if actual != testCase.expected {
			t.Error(testCase)
		}
	}
}
