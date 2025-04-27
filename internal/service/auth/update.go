package auth

import (
	"context"
	"github.com/Ippolid/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, id int64, info *model.UserInfo) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.authRepository.UpdateUser(ctx, id, *info)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
