package config

import (
	"os"

	"github.com/pkg/errors"
)

const (
	refreshTokenSecretKey = "REFRESH_TOKEN"
	accessTokenSecretKey  = "ACCESS_TOKEN"
)

// SwaggerConfig интерфейс для конфигурации Swagger-сервера
type JWTConfig interface {
	RefreshToken() string
	AccessToken() string
}

type jwtConfig struct {
	refreshkey string
	accesskey  string
}

// NewSwaggerConfig создает новую конфигурацию Swagger-сервера
func NewJWTConfig() (JWTConfig, error) {
	refreshkey := os.Getenv(refreshTokenSecretKey)
	if len(refreshkey) == 0 {
		return nil, errors.New("refreshkey not found")
	}

	accesskey := os.Getenv(accessTokenSecretKey)
	if len(accesskey) == 0 {
		return nil, errors.New("accesskey port not found")
	}

	return &jwtConfig{
		refreshkey: refreshkey,
		accesskey:  accesskey,
	}, nil
}

func (cfg *jwtConfig) RefreshToken() string {
	return cfg.refreshkey
}

func (cfg *jwtConfig) AccessToken() string {
	return cfg.accesskey
}
