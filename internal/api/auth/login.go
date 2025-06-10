package auth

import (
	"context"

	"github.com/Ippolid/auth/internal/converter"
	"github.com/Ippolid/auth/pkg/auth_v1"
)

// Login обрабатывает запрос на вход в систему
func (i *Controller) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	resp, err := i.authService.Login(ctx, *converter.ToLoginFromAuthAPI(req))
	if err != nil {
		return nil, err
	}

	return &auth_v1.LoginResponse{
		RefreshToken: resp.Token,
	}, nil
}
