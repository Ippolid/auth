package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const yamldir = "ACCESS_YAML_DIR"

// AccessConfig интерфейс для получения конфигурации доступа к эндпоинтам.
type AccessConfig interface {
	CFG() map[string]bool
}
type accessConfigImpl struct {
	Endpoints map[string]bool `yaml:"endpoints"`
}

// NewAccessConfig создает новую конфигурацию доступа из YAML-файла.
// Путь к директории YAML-файла берется из переменной окружения, указанной в константе yamldir.
// Имя файла по умолчанию "access.yaml".
func NewAccessConfig() (AccessConfig, error) {
	yamlPathDir := os.Getenv(yamldir) // Используем вашу константу yamldir
	if yamlPathDir == "" {
		return nil, errors.Errorf("переменная окружения %s не установлена", yamldir)
	}

	// Имя файла конфигурации можно сделать константой или параметром

	data, err := os.ReadFile(yamlPathDir) //nolint:gosec
	if err != nil {
		return nil, errors.Wrapf(err, "ошибка чтения файла конфигурации '%s'", yamlPathDir)
	}

	var cfg accessConfigImpl
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, errors.Wrapf(err, "ошибка парсинга YAML из '%s'", yamlPathDir)
	}

	if cfg.Endpoints == nil {
		// Инициализируем мапу, если YAML пуст или не содержит ключ 'endpoints'
		cfg.Endpoints = make(map[string]bool)
	}

	return &cfg, nil
}

func (cfg *accessConfigImpl) CFG() map[string]bool {
	// Возвращаем мапу эндпоинтов
	return cfg.Endpoints

}

// HasAccess проверяет, имеет ли указанная роль доступ к эндпоинту.
