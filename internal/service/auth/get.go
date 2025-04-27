package auth

import (
	"context"
	"fmt"
	"github.com/Ippolid/auth/internal/model"
	"time"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.authRepository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	err1 := s.authRepository.MakeLog(ctx, model.Log{
		Method:    "GET",
		CreatedAt: time.Now(),
		Ctx:       fmt.Sprintf("%v", ctx),
	})
	if err1 != nil {
		return nil, err1
	}

	return user, nil
}
