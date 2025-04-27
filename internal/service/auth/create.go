package auth

import (
	"context"
	"fmt"
	"github.com/Ippolid/auth/internal/model"
	"time"
)

func (s *serv) Create(ctx context.Context, info *model.User) (int64, error) {
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

	return id, nil
}
