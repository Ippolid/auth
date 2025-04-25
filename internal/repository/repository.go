package repository

import (
	"context"
	"github.com/Ippolid/auth/internal/model"
)

type AuthRepository interface {
	InsertUser(ctx context.Context, user model.User) (int64, error)
	GetUser(ctx context.Context, id int) (*model.User, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, id int, info model.UserInfo) error
}
