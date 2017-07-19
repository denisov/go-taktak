# go-taktak

Утилита для сбора статистики ответов сайта https://taktaktak.ru/

##### Особенности:
* Собирает отчёт для эксперта за указаннный месяц и год
* Формируются XLSX файлы отчёта
* Отчёт отправляется на email
* Можно задать маппинг соответствия id пользователя и его фамилии

##### Установка
* Создать рабочую директорию, по умолчанию `~/taktak`
* Скопировать в неё `config.yml`, внести в него соответствующие правки
* отчёты будут создаваться в директории `reports` внутри базовой директории 
  
##### Запуск
```bash
$ go-taktak --help                                                                                                              
Usage: run.php [--user=<user_id>] [--month=<month>] [--year=<year>] [--sleep=<sec>] [--month-deep=<n>] [--base-dir=<dir>]

Options:
  --user=<user_id>  id пользователя. Мила 15525. Марина 17881 [default: 15525]
  --month=<month>   месяц. По-умолчанию текущий.
  --year=<year>     год. По-умолчанию текущий.
  --month-deep=<n>  на сколько меяцев в прошлое уходить [default: 3]
  --sleep=<sec>     сон в сек между запросами [default: 1].
  --base-dir=<dir>  базовая директория, в ней ищется конфиг и создаются отчёты [default: ~/taktak]

```  
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6bb3bf5ca2c547e1990e76792bdd566d)](https://www.codacy.com/app/elgato.andrey/go-taktak?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=denisov/go-taktak&amp;utm_campaign=Badge_Grade)  
