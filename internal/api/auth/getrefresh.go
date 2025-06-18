package auth

import (
	"context"

	"github.com/Ippolid/auth/internal/converter"
	"github.com/Ippolid/auth/pkg/auth_v1"
)

// GetRefreshToken обрабатывает запрос на получение refresh-токена
func (i *Controller) GetRefreshToken(ctx context.Context, req *auth_v1.GetRefreshTokenRequest) (*auth_v1.GetRefreshTokenResponse, error) {
	resp, err := i.authService.GetRefreshToken(ctx, *converter.ToGetRefreshTokenFromDesc(req))
	if err != nil {
		return nil, err
	}

	return &auth_v1.GetRefreshTokenResponse{
		RefreshToken: resp.RefreshToken,
	}, nil
}
