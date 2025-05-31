package user

import (
	"github.com/Ippolid/auth/internal/service"
	"github.com/Ippolid/auth/pkg/user_v1"
)

// Controller реализует интерфейс AuthV1Server и AuthService
type Controller struct {
	user_v1.UnimplementedUserV1Server
	userService service.UserService
}

// NewController создает новый экземпляр Controller
func NewController(userService service.UserService) *Controller {
	return &Controller{
		userService: userService,
	}
}
