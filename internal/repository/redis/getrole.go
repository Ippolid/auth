package redis

import (
	"context"
	"fmt"

	"github.com/Ippolid/auth/internal/model"
)

func (c cache) GetRole(ctx context.Context, username string) (*bool, error) {
	result, err := c.cl.Get(ctx, username)
	if err != nil || result == nil {
		if result == nil {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("redis Get error: %w", err)
	}
	// Преобразуем результат в булево значение
	role, ok := result.(bool)
	if !ok {
		return nil, fmt.Errorf("unexpected role type in cache for user %s", username)
	}

	return &role, nil
}
