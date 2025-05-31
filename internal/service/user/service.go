package user

import (
	"github.com/Ippolid/auth/internal/repository"
	"github.com/Ippolid/auth/internal/service"
	"github.com/Ippolid/platform_libary/pkg/db"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
	cache          repository.CacheInterface
}

// NewService создает новый экземпляр AuthService
func NewService(
	userRepository repository.UserRepository,
	txManager db.TxManager,
	cache repository.CacheInterface,
) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
		cache:          cache,
	}
}

// NewMockService создает новый экземпляр AuthService для тестирования
func NewMockService(deps ...interface{}) service.UserService {
	srv := serv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.UserRepository:
			srv.userRepository = s
		case db.TxManager:
			srv.txManager = s
		}

	}

	return &srv
}
