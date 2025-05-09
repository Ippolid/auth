package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Ippolid/auth/internal/api/auth"
	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/service"
	"github.com/Ippolid/auth/internal/service/mocks"
	"github.com/Ippolid/auth/pkg/auth_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestController_Get(t *testing.T) {
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *auth_v1.GetRequest
	}

	var (
		ctx       = context.Background()
		id  int64 = 123

		serviceErr = fmt.Errorf("service error")

		req = &auth_v1.GetRequest{Id: id}

		user = &model.User{
			ID: id,
			User: model.UserInfo{
				Name:  "Test Name",
				Email: "test@example.com",
			},
			Role:      true, // true соответствует Role_ADMIN
			Password:  "password",
			CreatedAt: time.Time{}, // Нулевое значение времени
		}
	)

	tests := []struct {
		name            string
		args            args
		wantResp        *auth_v1.GetResponse
		wantErr         error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success",
			args: args{ctx: ctx, req: req},
			wantResp: &auth_v1.GetResponse{
				User: &auth_v1.UserGet{
					Id:        user.ID,
					Info:      &auth_v1.UserInfo{Name: user.User.Name, Email: user.User.Email},
					Role:      auth_v1.Role_ADMIN,           // Соответствует user.Role = true
					CreatedAt: timestamppb.New(time.Time{}), // Нулевое значение времени
					UpdatedAt: timestamppb.New(time.Time{}), // Нулевое значение времени
				},
			},
			wantErr: nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
		},
		{
			name:     "service error",
			args:     args{ctx: ctx, req: req},
			wantResp: nil,
			wantErr:  serviceErr,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			authService := tt.authServiceMock(mc)
			ctrl := auth.NewController(authService)
			resp, err := ctrl.Get(tt.args.ctx, tt.args.req)
			fmt.Printf("User id: %d\n", id) // Для отладки
			require.Equal(t, tt.wantResp, resp)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
