package user

import (
	"context"
	"github.com/Ippolid/auth/internal/converter"
	"github.com/Ippolid/auth/internal/logger"
	"github.com/Ippolid/auth/pkg/user_v1"
	"go.uber.org/zap"
)

// Get реализует метод получения пользователя по ID
func (i *Controller) Get(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	logger.Info("Get user request",
		zap.Int64("UserID", req.GetId()),
	)

	user1, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	responce := converter.ToUserAPIFromUserGet(user1)

	return responce, nil
}
