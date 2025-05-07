package repository

import (
	"context"

	"github.com/Ippolid/auth/internal/model"
)

// AuthRepository интерфейс для работы с репозиторием
type AuthRepository interface {
	CreateUser(ctx context.Context, user model.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	DeleteUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, id int64, info model.UserInfo) error
	MakeLog(ctx context.Context, log model.Log) error
}

// CacheInterface интерфейс для работы с кэшем
type CacheInterface interface {
	Create(ctx context.Context, id int64, user model.User) error
	Get(ctx context.Context, id int64) (*model.User, error)
}
