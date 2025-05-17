package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Ippolid/auth/internal/model"
)

//func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
//	var (
//		userProfile *model.User
//		errCache    error
//		err         error
//	)
//
//	userProfile, errCache = s.cache.Get(ctx, id)
//	fmt.Println(userProfile, errCache, errors.Is(errCache, model.ErrUserNotFound))
//	if errCache != nil {
//		if errors.Is(errCache, model.ErrUserNotFound) {
//			err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
//				var errTx error
//				userProfile, errTx = s.authRepository.GetUser(ctx, id)
//				if errTx != nil {
//					return fmt.Errorf("error getting user profile: %w", errTx)
//				}
//
//				fmt.Println(userProfile)
//
//				errTx = s.authRepository.MakeLog(ctx, model.Log{
//					Method:    "GET",
//					CreatedAt: time.Now(),
//					Ctx:       fmt.Sprintf("%v", ctx),
//				})
//				if errTx != nil {
//					return fmt.Errorf("error creating log: %w", errTx)
//				}
//
//				if errTx = s.cache.Create(ctx, userProfile.ID, *userProfile); errTx != nil {
//					return fmt.Errorf("error caching user profile: %w", errTx)
//				}
//
//				return nil
//			})
//			if err != nil {
//				fmt.Println(err)
//				return nil, err
//			}
//		}
//		return nil, fmt.Errorf("error with cache: %w", errCache)
//	}
//	fmt.Println(userProfile)
//	return userProfile, nil
//
//}

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	var (
		userProfile *model.User
		errCache    error
		err         error
	)

	userProfile, errCache = s.cache.Get(ctx, id)
	fmt.Println(userProfile, errCache, errors.Is(errCache, model.ErrUserNotFound))
	if errCache != nil {
		if errors.Is(errCache, model.ErrUserNotFound) {
			err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
				var errTx error
				userProfile, errTx = s.authRepository.GetUser(ctx, id)
				if errTx != nil {
					return fmt.Errorf("error getting user profile: %w", errTx)
				}

				fmt.Println(userProfile)

				errTx = s.authRepository.MakeLog(ctx, model.Log{
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
