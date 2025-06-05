package auth

import (
	"context"

	"github.com/Ippolid/auth/internal/converter"
	"github.com/Ippolid/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Check обрабатывает запрос на проверку доступа
func (i *Controller) Check(ctx context.Context, req *auth_v1.CheckRequest) (*emptypb.Empty, error) {
	err := i.authService.Check(ctx, *converter.ToCheckAccessFromDesc(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
