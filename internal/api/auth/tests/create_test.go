package auth_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Ippolid/auth/internal/service"
	"github.com/brianvoe/gofakeit/v6"

	"github.com/Ippolid/auth/internal/api/auth"
	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/service/mocks"
	"github.com/Ippolid/auth/pkg/auth_v1"
	desc "github.com/Ippolid/auth/pkg/auth_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestController_Create(t *testing.T) {
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

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
		wantResp        *auth_v1.CreateResponse
		wantErr         error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success",
			args: args{ctx: ctx, req: req},
			wantResp: &auth_v1.CreateResponse{
				Id: id,
			},
			wantErr: nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, &info).Return(id, nil)
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
				mock.CreateMock.Return(int64(0), serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			noteServiceMock := tt.authServiceMock(mc)
			api := auth.NewController(noteServiceMock)
			resp, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.wantResp, resp)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
