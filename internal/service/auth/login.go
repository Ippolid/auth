package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/utils"
	"time"
)

func (s *serv) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	var resp model.LoginResponse
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		user, errTx := s.authRepository.Login(ctx, req)
		if errTx != nil {
			return errTx
		}

		err1 := s.authRepository.MakeLog(ctx, model.Log{
			Method:    "Login",
			CreatedAt: time.Now(),
			Ctx:       fmt.Sprintf("%v", ctx),
		})

		if err1 != nil {
			return err1
		}

		refreshToken, err := utils.GenerateToken(*user,
			[]byte(s.token.RefreshToken()),
			refreshTokenExpiration,
		)

		if err != nil {
			return errors.New("failed to generate token")
		}

		resp.Token = refreshToken

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
