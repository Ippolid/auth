package auth

import (
	"context"
	"fmt"
	"github.com/Ippolid/auth/internal/model"
	"time"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.authRepository.DeleteUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		err := s.authRepository.MakeLog(ctx, model.Log{
			Method:    "Delete",
			CreatedAt: time.Now(),
			Ctx:       fmt.Sprintf("%v", ctx),
		})
		if err != nil {
			return err
		}

		_, errTx = s.authRepository.GetUser(ctx, id)
		if errTx == nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
