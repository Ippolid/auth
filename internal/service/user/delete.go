package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Ippolid/auth/internal/model"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.DeleteUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		err := s.userRepository.MakeLog(ctx, model.Log{
			Method:    "Delete",
			CreatedAt: time.Now(),
			Ctx:       fmt.Sprintf("%v", ctx),
		})
		if err != nil {
			return err
		}

		_, errTx = s.userRepository.GetUser(ctx, id)
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
