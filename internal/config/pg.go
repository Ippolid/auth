package config

import (
	"errors"
	"os"
)

const (
	// DSNKey - ключ для строки подключения к PostgreSQL
	DSNKey = "PG_DSN"
	//DSNlogger - ключ для строки подключения к Дб ЛОгера
	DSNlogger = "PG_DSN_LOGER"
)

// PGConfig представляет интерфейс для получения строки подключения к PostgreSQL
type PGConfig interface {
	DSN() string
	LogerDSN() string
}

type pgConfig struct {
	dsn      string
	logerDSN string
}

// NewPGConfig создает новый экземпляр PGConfig, получая DSN из переменной окружения
func NewPGConfig() (PGConfig, error) {
	dsn := os.Getenv(DSNKey)
	logerdsn := os.Getenv(DSNlogger)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}
	if len(logerdsn) == 0 {
		return nil, errors.New("pg logerdsn not found")
	}

	return &pgConfig{
		dsn:      dsn,
		logerDSN: logerdsn,
	}, nil
}

// возвратит DSN для подключения к PostgreSQL
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}

func (cfg *pgConfig) LogerDSN() string {
	return cfg.logerDSN
}
