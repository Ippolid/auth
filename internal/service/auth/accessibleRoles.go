package auth

import (
	"context"
	"fmt"
	"github.com/Ippolid/auth/internal/model"
	"time"
)

var accessibleRoles map[string]string

func (s *serv) accessibleRoles(ctx context.Context) (map[string]string, error) {
	fmt.Println(accessibleRoles)
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]string)

		// Пытаемся получить эндпоинты админа из кеша
		adminEndpoints, errCacheAdmin := s.cache.GetRoleEndpoints(ctx, true)
		fmt.Println(adminEndpoints, errCacheAdmin)
		if errCacheAdmin != nil {

			// Данных в кеше нет, выполняем загрузку из базы
			err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
				// Получаем эндпоинты для админа из базы
				endpoints, errTx := s.authRepository.GetUsersAccess(ctx, true)
				if errTx != nil {
					return fmt.Errorf("ошибка получения админских эндпоинтов из БД: %w", errTx)
				}

				fmt.Println(endpoints)

				// Записываем лог операции
				errLog := s.authRepository.MakeLog(ctx, model.Log{
					Method:    "Get admin endpoints",
					CreatedAt: time.Now(),
					Ctx:       fmt.Sprintf("%v", ctx),
				})
				if errLog != nil {
					return fmt.Errorf("ошибка создания лога: %w", errLog)
				}

				// Сохраняем в кеш для будущих запросов
				errCache := s.cache.CreateRoleEndpoints(ctx, true, endpoints)
				if errCache != nil {
					return fmt.Errorf("ошибка кеширования админских эндпоинтов: %w", errCache)
				}

				adminEndpoints = endpoints
				return nil
			})
			if err != nil {
				return nil, err
			}
		}

		// Заполняем мапу для эндпоинтов админа
		for _, endpoint := range adminEndpoints {
			accessibleRoles[endpoint] = "admin"
		}
	}

	return accessibleRoles, nil
}
