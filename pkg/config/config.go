package config

import (
	"github.com/joho/godotenv"
	"newsaggr/pkg/logger"
	"os"
	"path/filepath"
	"strconv"
)

var Cfg *Config

// Config - конфигурация
type Config struct {
	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	DatabasePassword string
	Database         string
	RunPort          int
}

const debug = true

// Init - инициализация переменных окружения
func Init() error {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Error("Ошибка при попытке получить путь к директории: %v", err)
		return err
	}

	if debug {
		pwd = "C:\\Users\\Alex\\GolandProjects\\newsaggr"
	}

	pathToEnv := filepath.Join(pwd, ".env")

	err = godotenv.Load(pathToEnv)
	if err != nil {
		logger.Error("Ошибка при загрузке .env файла: %v", err)
		return err
	}

	DatabasePort, _ := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	RunPort, _ := strconv.Atoi(os.Getenv("RUN_PORT"))

	Cfg = &Config{
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     DatabasePort,
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		Database:         os.Getenv("DATABASE"),
		RunPort:          RunPort,
	}

	return nil
}
