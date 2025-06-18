package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Ippolid/auth/pkg/user_v1"

	"github.com/Ippolid/auth/internal/api/user"
	"github.com/Ippolid/auth/internal/service"
	"github.com/Ippolid/auth/internal/service/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestController_Delete(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *user_v1.DeleteRequest
	}

	var (
		ctx          = context.Background()
		id     int64 = 123
		req          = &user_v1.DeleteRequest{Id: id}
		svcErr       = errors.New("service error")
	)

	tests := []struct {
		name            string
		args            args
		wantResp        *emptypb.Empty
		wantErr         error
		userServiceMock userServiceMockFunc
	}{
		{
			name:     "success",
			args:     args{ctx: ctx, req: req},
			wantResp: &emptypb.Empty{},
			wantErr:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name:     "service error",
			args:     args{ctx: ctx, req: req},
			wantResp: nil,
			wantErr:  svcErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(svcErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			userService := tt.userServiceMock(mc)
			ctrl := user.NewController(userService)
			resp, err := ctrl.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.wantResp, resp)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, resp)
			}
		})
	}
}
