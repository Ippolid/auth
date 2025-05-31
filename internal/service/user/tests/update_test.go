package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/repository"
	repoMocks "github.com/Ippolid/auth/internal/repository/mocks"
	"github.com/Ippolid/auth/internal/service/mocks"
	"github.com/Ippolid/auth/internal/service/user"
	"github.com/Ippolid/platform_libary/pkg/db"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx  context.Context
		id   int64
		info model.UserInfo
	}

	var (
		ctx     = context.Background()
		id      = gofakeit.Int64()
		repoErr = fmt.Errorf("repo error")
		info    = model.UserInfo{
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
		}
	)

	tests := []struct {
		name               string
		args               args
		wantErr            error
		userRepositoryMock userRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				id:   id,
				info: info,
			},
			wantErr: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, id, info).Return(nil)
				mock.MakeLogMock.Set(func(_ context.Context, log model.Log) error {
					if log.Method != "Update" || log.Ctx != Ctxstring {
						return fmt.Errorf("unexpected log entry: %+v", log)
					}
					return nil
				})
				return mock

			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "UpdateUser error case",
			args: args{
				ctx:  ctx,
				id:   id,
				info: info,
			},
			wantErr: repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, id, info).Return(repoErr)
				// MakeLog и GetUser не должны вызываться при ошибке UpdateUser
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			t.Cleanup(mc.Finish)

			userRepoMock := tt.userRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)
			service := user.NewService(userRepoMock, txManagerMock, nil)

			err := service.Update(tt.args.ctx, tt.args.id, &tt.args.info)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
