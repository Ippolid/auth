package user

import (
	"context"
	"github.com/Ippolid/auth/internal/converter"
	"github.com/Ippolid/auth/internal/logger"
	"github.com/Ippolid/auth/pkg/user_v1"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// Get реализует метод получения пользователя по ID
func (i *Controller) Get(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	logger.Info("Get user request",
		zap.Int64("UserID", req.GetId()),
	)

	span, ctx := opentracing.StartSpanFromContext(ctx, "get note")
	defer span.Finish()

	span.SetTag("id", req.GetId())

	user1, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	responce := converter.ToUserAPIFromUserGet(user1)

	return responce, nil
}
