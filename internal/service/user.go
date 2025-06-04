package service

import (
	"context"

	"github.com/Ippolid/auth/internal/model"
)

// UserService интерфейс для работы с пользователями
type UserService interface {
	Create(ctx context.Context, info *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, info *model.UserInfo) error
}
