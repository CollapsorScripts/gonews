package rss

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const test = true

type Data struct {
	URLS          []string `json:"rss"`
	RequestPeriod uint     `json:"request_period"`
}

// GetData - парсит конфигурацию и возвращает ссылку на структуру с данными
func GetData() (*Data, error) {
	data := &Data{}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if test {
		pwd = "C:\\Users\\Alex\\GolandProjects\\newsaggr"
	}

	pathToConfig := filepath.Join(pwd, "config.json")

	cfg, err := os.ReadFile(pathToConfig)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(cfg, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
