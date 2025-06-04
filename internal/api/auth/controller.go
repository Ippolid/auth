package auth

import (
	"github.com/Ippolid/auth/internal/service"
	"github.com/Ippolid/auth/pkg/auth_v1"
)

type Controller struct {
	auth_v1.UnimplementedAuthServer
	authService service.AuthService
}

// NewController создает новый экземпляр Controller
func NewController(authService service.AuthService) *Controller {
	return &Controller{
		authService: authService,
	}
}
