package user_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Ippolid/auth/internal/service"
	"github.com/brianvoe/gofakeit/v6"

	"github.com/Ippolid/auth/internal/api/user"
	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/service/mocks"
	"github.com/Ippolid/auth/pkg/user_v1"
	desc "github.com/Ippolid/auth/pkg/user_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestController_Create(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()

		id = gofakeit.Int64()

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Info: &desc.UserInfoCreate{
				User: &desc.UserInfo{
					Name:  gofakeit.Name(),
					Email: gofakeit.Email(),
				},
				Password:        "password",
				PasswordConfirm: "password",
				Role:            desc.Role_ADMIN,
			},
		}

		info = model.User{
			User: model.UserInfo{
				Name:  req.GetInfo().GetUser().GetName(),
				Email: req.GetInfo().GetUser().GetEmail(),
			},
			Password: req.GetInfo().GetPassword(),
			Role:     req.GetInfo().GetRole() > 0,
		}
	)

	tests := []struct {
		name            string
		args            args
		wantResp        *user_v1.CreateResponse
		wantErr         error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success",
			args: args{ctx: ctx, req: req},
			wantResp: &user_v1.CreateResponse{
				Id: id,
			},
			wantErr: nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, &info).Return(id, nil)
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
				mock.CreateMock.Return(int64(0), serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			noteServiceMock := tt.userServiceMock(mc)
			api := user.NewController(noteServiceMock)
			resp, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.wantResp, resp)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
