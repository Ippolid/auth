package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Ippolid/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, id int64, info *model.UserInfo) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.authRepository.UpdateUser(ctx, id, *info)
		if errTx != nil {
			return errTx
		}

		err := s.authRepository.MakeLog(ctx, model.Log{
			Method:    "Update",
			CreatedAt: time.Now(),
			Ctx:       fmt.Sprintf("%v", ctx),
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
