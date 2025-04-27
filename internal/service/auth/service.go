package auth

import (
	"github.com/Ippolid/auth/internal/client/db"
	"github.com/Ippolid/auth/internal/repository"
	"github.com/Ippolid/auth/internal/service"
)

type serv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
}

func NewService(
	authRepository repository.AuthRepository,
	txManager db.TxManager,
) service.AuthService {
	return &serv{
		authRepository: authRepository,
		txManager:      txManager,
	}
}
