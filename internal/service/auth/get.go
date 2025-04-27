package auth

import (
	"context"
	"github.com/Ippolid/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.authRepository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
