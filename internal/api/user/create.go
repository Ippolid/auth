package user

import (
	"context"

	"github.com/Ippolid/auth/internal/converter"
	"github.com/Ippolid/auth/pkg/user_v1"

	"log"
)

// Create реализует метод создания пользователя
func (i *Controller) Create(ctx context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToUserCreateFromUserAPI(req))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted note with id: %d", id)

	return &user_v1.CreateResponse{
		Id: id,
	}, nil
}
