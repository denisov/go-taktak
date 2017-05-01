package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/docopt/docopt-go"
	"os/user"
	"path/filepath"
	"strings"
)

// проблема на которую отвечал эксперт
type problem struct {
	title        string
	url          string
	problemDate  time.Time
	solutionDate time.Time
}

var cfg config

func main() {
	log.SetOutput(os.Stdout)

	usage := `
Usage: run.php [--user=<user_id>] [--month=<month>] [--year=<year>] [--sleep=<sec>] [--month-deep=<n>] [--base-dir=<dir>]

Options:
  --user=<user_id>  id пользователя. Мила 15525. Марина 17881 [default: 15525]
  --month=<month>   месяц. По-умолчанию текущий.
  --year=<year>     год. По-умолчанию текущий.
  --month-deep=<n>  на сколько меяцев в прошлое уходить [default: 3]
  --sleep=<sec>     сон в сек между запросами [default: 1].
  --base-dir=<dir>  базовая директория, в ней ищется конфиг и создаются отчёты [default: ~/taktak]
`

	arguments, err := docopt.Parse(usage, nil, true, "1.0.0", false)
	if err != nil {
		log.Fatal("Не могу распарсить параметры")
	}

	var month int
	if arguments["--month"] == nil {
		month = int(time.Now().Month())
	} else {
		month, err = strconv.Atoi(arguments["--month"].(string))
		if err != nil {
			log.Fatal(err)
		}
	}

	var year int
	if arguments["--year"] == nil {
		year = time.Now().Year()
	} else {
		year, err = strconv.Atoi(arguments["--year"].(string))
		if err != nil {
			log.Fatal(err)
		}
	}

	userId, err := strconv.Atoi(arguments["--user"].(string))
	if err != nil {
		log.Fatal(err)
	}

	monthDeep, err := strconv.Atoi(arguments["--month-deep"].(string))
	if err != nil {
		log.Fatal(err)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	baseDir := filepath.Clean(arguments["--base-dir"].(string))
	baseDir = strings.Replace(baseDir, "~", usr.HomeDir, 1)
	log.Printf("Базовая директория: %s", baseDir)

	cfg, err = newConfig(filepath.Join(baseDir, "config.yml"))
	if err != nil {
		log.Fatalln(err)
	}

	reportsDir := filepath.Join(baseDir, "reports")
	err = ensureDir(reportsDir)
	if err != nil {
		log.Fatalln(err)
	}

	problems := getUserStat(userId, month, year, monthDeep)

	reportName := getReportName(userId, month, year)
	reportPath := filepath.Join(reportsDir, reportName)
	writeReportFile(problems, reportPath)
	sendReport(reportPath, reportName)
}
