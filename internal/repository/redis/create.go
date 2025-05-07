package redis

import (
	"context"
	"fmt"
	"github.com/Ippolid/auth/internal/model"
	"strconv"
	"time"
)

func (c cache) Create(ctx context.Context, id int64, user model.User) error {
	idFormatted := strconv.FormatInt(id, 10)

	redisUser := toRedisModels(id, user)
	if err := c.cl.HashSet(ctx, idFormatted, redisUser); err != nil {
		return fmt.Errorf("failed to hash user: %w", err)
	}

	if err := c.cl.Expire(ctx, idFormatted, 5*time.Minute); err != nil {
		return fmt.Errorf("failed to set expiration for user: %w", err)
	}
	return nil
}
