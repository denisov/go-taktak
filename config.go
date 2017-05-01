package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type config struct {
	Users map[int]string
	Email struct {
		From     string
		To       string
		Password string
		Host     string
		Port     string
	}
}

func newConfig(fileName string) (config, error) {
	newConfig := config{}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return newConfig, fmt.Errorf("Не могу найти файл конфига: %s", fileName)
	}

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return newConfig, fmt.Errorf("Чтение файла %s провалилось: %v", fileName, err)
	}

	if err := yaml.Unmarshal(file, &newConfig); err != nil {
		return newConfig, fmt.Errorf("unmarshal failed: %v", err)
	}
	return newConfig, nil
}
