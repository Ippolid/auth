package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Ippolid/auth/internal/model"

	redismodels "github.com/Ippolid/auth/internal/repository/model"
)

const customTimeFormat = "2006-01-02 15:04:05.999999999"

func toRedisModels(id int64, user model.User) redismodels.UserRedis {
	idStr, timeNow := strconv.FormatInt(id, 10), time.Now()

	return redismodels.UserRedis{
		ID:        idStr,
		Name:      *user.User.Name,
		Email:     *user.User.Email,
		Password:  user.Password,
		Role:      strconv.FormatBool(user.Role),
		CreatedAt: timeNow.Format(customTimeFormat),
	}
}

func toServiceModels(user redismodels.UserRedis) (*model.User, error) {
	createdAt, err := time.Parse(customTimeFormat, user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error with time parse CreatedAt: %w", err)
	}

	id, err := strconv.ParseInt(user.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error with parse ID: %w", err)
	}

	role, err := strconv.ParseBool(user.Role)
	if err != nil {
		return nil, fmt.Errorf("error with parse Role: %w", err)
	}

	user1 := model.UserInfo{
		Name:  &user.Name,
		Email: &user.Email,
	}
	return &model.User{
		ID:        id,
		User:      user1,
		Password:  user.Password,
		Role:      role,
		CreatedAt: createdAt,
	}, nil
}
