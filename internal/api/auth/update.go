package auth

import (
	"context"
	"log"

	"github.com/Ippolid/auth/internal/converter"
	"github.com/Ippolid/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Update реализует метод обновления пользователя
func (i *Implementation) Update(ctx context.Context, req *auth_v1.UpdateRequest) (*emptypb.Empty, error) {
	user := converter.ToUserInfoFromService(req)
	err := i.authService.Update(ctx, req.GetId(), user)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &emptypb.Empty{}, nil
}
