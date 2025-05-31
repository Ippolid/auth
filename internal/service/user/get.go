package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Ippolid/auth/internal/model"
)

// Get получает профиль пользователя по ID
func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	var (
		userProfile *model.User
		errCache    error
		err         error
	)

	userProfile, errCache = s.cache.Get(ctx, id)
	if errCache != nil {
		if errors.Is(errCache, model.ErrUserNotFound) {
			err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
				var errTx error
				userProfile, errTx = s.userRepository.GetUser(ctx, id)
				if errTx != nil {
					return fmt.Errorf("error getting user profile: %w", errTx)
				}

				errTx = s.userRepository.MakeLog(ctx, model.Log{
					Method:    "GET",
					CreatedAt: time.Now(),
					Ctx:       fmt.Sprintf("%v", ctx),
				})
				if errTx != nil {
					return fmt.Errorf("error creating log: %w", errTx)
				}

				if errTx = s.cache.Create(ctx, userProfile.ID, *userProfile); errTx != nil {
					return fmt.Errorf("error caching user profile: %w", errTx)
				}

				return nil
			})
			if err != nil {
				return nil, err
			}
			// если успешно получили из БД, возвращаем результат
			return userProfile, nil
		}
		// если ошибка другая — возвращаем ошибку кеша
		return nil, fmt.Errorf("error with cache: %w", errCache)
	}
	return userProfile, nil
}
