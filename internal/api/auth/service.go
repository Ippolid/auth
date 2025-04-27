package auth

import (
	"github.com/Ippolid/auth/internal/service"
	"github.com/Ippolid/auth/pkg/auth_v1"
)

type Implementation struct {
	auth_v1.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
