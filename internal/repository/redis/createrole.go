package redis

import (
	"context"
	"fmt"
	"time"
)

func (c cache) CreateRole(ctx context.Context, username string, role bool) error {
	// Используем команду Set, чтобы сохранить роль
	if err := c.cl.Set(ctx, username, role); err != nil {
		return fmt.Errorf("failed to set role for username %s: %w", username, err)
	}

	// Отдельно устанавливаем время жизни ключа
	if err := c.cl.Expire(ctx, username, 6*time.Minute); err != nil {
		return fmt.Errorf("failed to set expiration for username %s: %w", username, err)
	}

	return nil
}
