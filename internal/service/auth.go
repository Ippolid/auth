package service

import (
	"context"

	"github.com/Ippolid/auth/internal/model"
)

// AuthService интерфейс для работы с авторизацией
type AuthService interface {
	Check(ctx context.Context, request model.CheckRequest) error
	Login(ctx context.Context, request model.LoginRequest) (*model.LoginResponse, error)
	GetRefreshToken(ctx context.Context, request model.GetRefreshTokenRequest) (*model.GetRefreshTokenResponse, error)
	GetAccessToken(ctx context.Context, req model.GetAccessTokenRequest) (*model.GetAccessTokenResponse, error)
}
