package config

import (
	"github.com/joho/godotenv"
)

// Load загружает переменные окружения из файла .env
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
