package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/Ippolid/auth/internal/api/auth"
	"github.com/Ippolid/auth/internal/service"
	"github.com/Ippolid/auth/internal/service/mocks"
	"github.com/Ippolid/auth/pkg/auth_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestController_Delete(t *testing.T) {
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *auth_v1.DeleteRequest
	}

	var (
		ctx          = context.Background()
		id     int64 = 123
		req          = &auth_v1.DeleteRequest{Id: id}
		svcErr       = errors.New("service error")
	)

	tests := []struct {
		name            string
		args            args
		wantResp        *emptypb.Empty
		wantErr         error
		authServiceMock authServiceMockFunc
	}{
		{
			name:     "success",
			args:     args{ctx: ctx, req: req},
			wantResp: &emptypb.Empty{},
			wantErr:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name:     "service error",
			args:     args{ctx: ctx, req: req},
			wantResp: nil,
			wantErr:  svcErr,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(svcErr)
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
