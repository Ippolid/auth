package repository

import (
	"context"

	"github.com/Ippolid/auth/internal/model"
)

// UserRepository интерфейс для работы с репозиторием
type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	DeleteUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, id int64, info model.UserInfo) error
	MakeLog(ctx context.Context, log model.Log) error
}

type AuthRepository interface {
	Login(ctx context.Context, user model.LoginRequest) (*model.UserInfoJwt, error)
	MakeLog(ctx context.Context, log model.Log) error
	GetUserRole(ctx context.Context, username string) (bool, error)
	GetUsersAccess(ctx context.Context, isAdmin bool) ([]string, error)
}

// CacheInterface интерфейс для работы с кэшем
type CacheInterface interface {
	Create(ctx context.Context, id int64, user model.User) error
	Get(ctx context.Context, id int64) (*model.User, error)
	GetRole(ctx context.Context, username string) (*bool, error)
	CreateRole(ctx context.Context, username string, role bool) error
	CreateRoleEndpoints(ctx context.Context, isAdmin bool, endpoints []string) error
	GetRoleEndpoints(ctx context.Context, isAdmin bool) ([]string, error)
}
