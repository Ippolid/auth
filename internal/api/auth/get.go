package auth

import (
	"context"
	"log"

	"github.com/Ippolid/auth/internal/converter"
	"github.com/Ippolid/auth/pkg/auth_v1"
)

// Get реализует метод получения пользователя по ID
func (i *Implementation) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	user1, err := i.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	responce := converter.ToDescFromAuthGet(user1)

	return responce, nil
}
