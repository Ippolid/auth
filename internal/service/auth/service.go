package auth

import (
	"github.com/Ippolid/auth/internal/repository"
	"github.com/Ippolid/auth/internal/service"
	"github.com/Ippolid/platform_libary/pkg/db"
)

type serv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
}

// NewService создает новый экземпляр AuthService
func NewService(
	authRepository repository.AuthRepository,
	txManager db.TxManager,
) service.AuthService {
	return &serv{
		authRepository: authRepository,
		txManager:      txManager,
	}
}

// NewMockService создает новый экземпляр AuthService для тестирования
func NewMockService(deps ...interface{}) service.AuthService {
	srv := serv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.AuthRepository:
			srv.authRepository = s
		case db.TxManager:
			srv.txManager = s
		}

	}

	return &srv
}
