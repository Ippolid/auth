package auth

import "context"

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.authRepository.DeleteUser(ctx, id)
		if errTx != nil {
			return errTx
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
