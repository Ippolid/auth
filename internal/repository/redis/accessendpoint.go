package redis

import (
	"context"
	"fmt"
	"github.com/Ippolid/auth/internal/model"
	"github.com/gomodule/redigo/redis"
)

// CreateRoleEndpoints сохраняет список эндпоинтов для указанной роли
func (c *cache) CreateRoleEndpoints(ctx context.Context, isAdmin bool, endpoints []string) error {
	if len(endpoints) == 0 {
		return nil
	}

	// Формируем ключ в формате "role:admin" или "role:user"
	roleKey := "role:user"
	if isAdmin {
		roleKey = "role:admin"
	}

	// Преобразуем слайс строк в массив интерфейсов для Redis
	args := make([]interface{}, len(endpoints)+1)
	args[0] = roleKey
	for i, endpoint := range endpoints {
		args[i+1] = endpoint
	}

	err := c.cl.Execute(ctx, func(_ context.Context, conn redis.Conn) error {
		// Удаляем старые значения (если есть)
		_, err := conn.Do("DEL", roleKey)
		if err != nil {
			return fmt.Errorf("failed to delete old endpoints: %w", err)
		}

		// Если есть эндпоинты для сохранения
		if len(endpoints) > 0 {
			// Сохраняем новые эндпоинты как элементы множества
			_, err = conn.Do("SADD", args...)
			if err != nil {
				return fmt.Errorf("failed to add endpoints to set: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to save endpoints for role (admin=%t): %w", isAdmin, err)
	}

	return nil
}

// GetRoleEndpoints получает список эндпоинтов для указанной роли
func (c *cache) GetRoleEndpoints(ctx context.Context, isAdmin bool) ([]string, error) {
	roleKey := "role:user"
	if isAdmin {
		roleKey = "role:admin"
	}

	var endpoints []string

	err := c.cl.Execute(ctx, func(_ context.Context, conn redis.Conn) error {
		values, err := redis.Strings(conn.Do("SMEMBERS", roleKey))
		if err != nil {
			if err == redis.ErrNil {
				return nil // Пустой результат, возвращаем пустой слайс
			}
			return err
		}

		endpoints = values
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints for role (admin=%t): %w", isAdmin, err)
	}

	if len(endpoints) == 0 {
		return nil, model.ErrUserNotFound
	}

	return endpoints, nil
}
