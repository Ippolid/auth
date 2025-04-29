package auth

import (
	"context"

	"github.com/Ippolid/auth/internal/converter"
	"github.com/Ippolid/auth/pkg/auth_v1"

	"log"
)

// Create реализует метод создания пользователя
func (i *Controller) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	id, err := i.authService.Create(ctx, converter.ToAuthCreateFromDesc(req))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted note with id: %d", id)

	return &auth_v1.CreateResponse{
		Id: id,
	}, nil
}
