package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ippolid/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.User) (int64, error) {
	if info == nil {
		return 0, fmt.Errorf("user info is nil")
	}

	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.authRepository.CreateUser(ctx, *info)
		if errTx != nil {
			return errTx
		}

		err := s.authRepository.MakeLog(ctx, model.Log{
			Method:    "Create",
			CreatedAt: time.Now(),
			Ctx:       fmt.Sprintf("%v", ctx),
		})
		if err != nil {
			return err
		}
		_, errTx = s.authRepository.GetUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	if s.cache != nil {
		if cacheErr := s.cache.Create(ctx, id, *info); cacheErr != nil {
			log.Printf("failed to cache user: %v", cacheErr)
		}
	} else {
		log.Println("cache is not initialized, skipping cache creation")
	}

	return id, nil
}
