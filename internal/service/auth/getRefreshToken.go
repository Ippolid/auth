package auth

import (
	"context"

	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) GetRefreshToken(_ context.Context, req model.GetRefreshTokenRequest) (*model.GetRefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(req.OldToken, []byte(s.token.RefreshToken()))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	refreshToken, err := utils.GenerateToken(model.UserInfoJwt{
		Username: claims.Username,
		Role:     claims.Role == "1" || claims.Role == "true",
	},
		[]byte(s.token.RefreshToken()),
		refreshTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &model.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}
