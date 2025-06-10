package auth

import (
	"context"

	"github.com/Ippolid/auth/internal/converter"
	"github.com/Ippolid/auth/pkg/auth_v1"
)

// GetAccessToken обрабатывает запрос на получение access-токена
func (i *Controller) GetAccessToken(ctx context.Context, req *auth_v1.GetAccessTokenRequest) (*auth_v1.GetAccessTokenResponse, error) {
	resp, err := i.authService.GetAccessToken(ctx, *converter.ToGetAccessTokenFromAuthApi(req))
	if err != nil {
		return nil, err
	}

	return &auth_v1.GetAccessTokenResponse{
		AccessToken: resp.AccessToken,
	}, nil
}
