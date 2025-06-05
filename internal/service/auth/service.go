package auth

import (
	"time"

	"github.com/Ippolid/auth/internal/config"
	"github.com/Ippolid/auth/internal/repository"
	"github.com/Ippolid/auth/internal/service"
	"github.com/Ippolid/platform_libary/pkg/db"
)

const (
	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 5 * time.Minute
)

type serv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
	cache          repository.CacheInterface
	token          config.JWTConfig
}

// NewService создает новый экземпляр AuthService
func NewService(
	authRepository repository.AuthRepository,
	txManager db.TxManager,
	cache repository.CacheInterface,
	token config.JWTConfig,
) service.AuthService {
	return &serv{
		authRepository: authRepository,
		txManager:      txManager,
		cache:          cache,
		token:          token,
	}
}
