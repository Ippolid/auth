package user

import (
	"context"
	"log"

	"github.com/Ippolid/auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete реализует метод удаления пользователя
func (i *Controller) Delete(ctx context.Context, req *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())

	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &emptypb.Empty{}, nil
}
