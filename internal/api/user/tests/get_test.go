package user_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	user1 "github.com/Ippolid/auth/internal/api/user"
	"github.com/Ippolid/auth/internal/logger"
	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/service"
	"github.com/Ippolid/auth/internal/service/mocks"
	"github.com/Ippolid/auth/pkg/user_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ptr[T any](v T) *T {
	return &v
}

func TestController_Get(t *testing.T) {
	// Инициализируем логгер для тестов
	logger.InitLocalLogger("Info")

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *user_v1.GetRequest
	}

	var (
		ctx       = context.Background()
		id  int64 = 123

		serviceErr = fmt.Errorf("service error")

		req = &user_v1.GetRequest{Id: id}

		user = &model.User{
			ID: id,
			User: model.UserInfo{
				Name:  ptr("Test Name"),
				Email: ptr("test@example.com"),
			},
			Role:      true, // true соответствует Role_ADMIN
			Password:  "password",
			CreatedAt: time.Time{}, // Нулевое значение времени
		}
	)

	tests := []struct {
		name            string
		args            args
		wantResp        *user_v1.GetResponse
		wantErr         error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success",
			args: args{ctx: ctx, req: req},
			wantResp: &user_v1.GetResponse{
				User: &user_v1.UserGet{
					Id:        user.ID,
					Info:      &user_v1.UserInfo{Name: *user.User.Name, Email: *user.User.Email},
					Role:      user_v1.Role_ADMIN,
					CreatedAt: timestamppb.New(time.Time{}),
					UpdatedAt: timestamppb.New(time.Time{}),
				},
			},
			wantErr: nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				// Используем minimock.AnyContext вместо ctx
				mock.GetMock.Expect(minimock.AnyContext, id).Return(user, nil)
				return mock
			},
		},
		{
			name:     "service error",
			args:     args{ctx: ctx, req: req},
			wantResp: nil,
			wantErr:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				// Используем minimock.AnyContext вместо ctx
				mock.GetMock.Expect(minimock.AnyContext, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			userService := tt.userServiceMock(mc)
			ctrl := user1.NewController(userService)
			resp, err := ctrl.Get(tt.args.ctx, tt.args.req)

			// Сравниваем ошибки и ответы
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, resp)
			}
		})
	}
}
