package auth

import (
	"github.com/Ippolid/auth/internal/service"
	"github.com/Ippolid/auth/pkg/auth_v1"
)

// Implementation реализует интерфейс AuthV1Server и AuthService
type Implementation struct {
	auth_v1.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation создает новый экземпляр Implementation
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
