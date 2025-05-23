package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Ippolid/auth/internal/model"
	redismodels "github.com/Ippolid/auth/internal/repository/model"
	"github.com/gomodule/redigo/redis"
)

func (c cache) Get(ctx context.Context, id int64) (*model.User, error) {
	var key string
	key = strconv.FormatInt(id, 10)

	userCache, err := c.cl.HGetAll(ctx, key)
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("redis HGetAll error: %w", err)
	}

	if len(userCache) == 0 {
		return nil, model.ErrUserNotFound
	}

	var userProfile redismodels.UserRedis
	err = redis.ScanStruct(userCache, &userProfile)
	if err != nil {
		return nil, fmt.Errorf("error scanning user profile: %w", err)
	}

	fmt.Printf("%v", userProfile)
	user, err := toServiceModels(userProfile)
	if err != nil {
		return nil, fmt.Errorf("error converting user profile: %w", err)
	}

	return user, nil
}
