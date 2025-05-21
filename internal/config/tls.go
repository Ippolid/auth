package config

import (
	"os"

	"github.com/pkg/errors"
)

const (
	tlsCertEnvName = "TLS_CERT_PATH"
	tlsKeyEnvName  = "TLS_KEY_PATH"
)

// TLSConfig интерфейс для получения путей до TLS-сертификата и ключа
type TLSConfig interface {
	CertFile() string
	KeyFile() string
}

// tlsConfig — внутренняя реализация TLSConfig
type tlsConfig struct {
	certPath string
	keyPath  string
}

// NewTLSConfig читает из окружения пути до .crt и .key и возвращает TLSConfig
func NewTLSConfig() (TLSConfig, error) {
	cert := os.Getenv(tlsCertEnvName)
	if cert == "" {
		return nil, errors.Errorf("%s not set", tlsCertEnvName)
	}

	key := os.Getenv(tlsKeyEnvName)
	if key == "" {
		return nil, errors.Errorf("%s not set", tlsKeyEnvName)
	}

	return &tlsConfig{
		certPath: cert,
		keyPath:  key,
	}, nil
}

// CertFile возвращает путь до TLS-сертификата (.crt)
func (cfg *tlsConfig) CertFile() string {
	return cfg.certPath
}

// KeyFile возвращает путь до приватного ключа (.key)
func (cfg *tlsConfig) KeyFile() string {
	return cfg.keyPath
}
