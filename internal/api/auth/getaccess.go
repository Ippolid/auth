package auth

import (
	"context"
	"github.com/Ippolid/auth/internal/converter"
	"github.com/Ippolid/auth/pkg/auth_v1"
)

func (i *Controller) GetAccessToken(ctx context.Context, req *auth_v1.GetAccessTokenRequest) (*auth_v1.GetAccessTokenResponse, error) {
	resp, err := i.authService.GetAccessToken(ctx, *converter.ToGetAccessTokenFromDesc(req))
	if err != nil {
		return nil, err
	}

	return &auth_v1.GetAccessTokenResponse{
		AccessToken: resp.AccessToken,
	}, nil
}
